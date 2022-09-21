package zaputils

import (
	"fmt"
	"runtime"
	"time"

	syslog "github.com/hashicorp/go-syslog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	logFile = "logs.log"
)

// NewProdLogger creates a new zap logger that is production ready.
// It consists of three "sub" loggers: a console logger, a file logger and a system logger.
// The configs used for each of these loggers are, respectively: ConsoleConfig, FileConfig and SysConfig.
// But, the levels of each config is overwritten as follows: the console logger will log Info level and above,
// the file and the system loggers will log Warn level and above.
//
// NewProdLogger does not implement a rotating logger for the file logger. If you want that functionality
// use NewRotatingProdLogger instead.
//
// The system logger only works on non windows systems.
func NewProdLogger(name, facility string, priority syslog.Priority) (*zap.Logger, error) {
	return NewCustomLogger(name, facility, priority, ConsoleConfig(), FileConfig(), SysConfig())
}

// NewRotatingProdLogger creates a new zap logger that is production ready.
// It consists of three "sub" loggers: a console logger, a file logger and a system logger.
// The configs used for each of these loggers are, respectively: ConsoleConfig, FileConfig and SysConfig.
// But, the levels of each config is overwritten as follows: the console logger will log Info level and above,
// the file and the system loggers will log Warn level and above.
//
// NewRotatingProdLogger implements a rotating logger system for the file logger.
// The max file size is 100mb, max age is 14 days, and rotated files will be compressed.
// If a custom configuration for the rotation is needed, call NewRotatingCustomLogger instead.
//
// The system logger only works on non windows systems.
func NewRotatingProdLogger(name, facility string, priority syslog.Priority) (*zap.Logger, error) {
	lj := &lumberjack.Logger{
		Filename: logFile,
		MaxSize:  100,
		MaxAge:   14,
		Compress: true,
	}
	return NewRotatingCustomLogger(name, facility, priority, ConsoleConfig(), FileConfig(), SysConfig(), lj)
}

// NewCustomLogger creates a new zap logger that consists of 3 sub loggers.
//
// cConfig is the config intended to be used as a console logger.
// The error sink will be grabbed from this config ErrorOutputPaths.
//
// fConfig is the config intended to be used as a file logger. The file path is `logs.log`.
// This logger does not implement a rotating logger, if you want that functonality call
// NewRotatingCustomLogger instead.
//
// sConfig is the config inteded to be used as a system logger.
// This will only work on non windows systems.
func NewCustomLogger(name, facility string, priority syslog.Priority, cConfig, fConfig, sConfig zap.Config) (*zap.Logger, error) {
	return newCustomLogger(name, facility, priority, cConfig, fConfig, sConfig, nil)
}

// NewRotatingCustomLogger creates a new zap logger that consists of 3 sub loggers.
//
// cConfig is the config intended to be used as a console logger.
// The error sink will be grabbed from this config ErrorOutputPaths.
//
// fConfig is the config intended to be used as a file logger. The file path is `logs.log`.
// This logger will implement a rotating log system based on lj.
//
// sConfig is the config inteded to be used as a system logger.
// This will only work on non windows systems.
func NewRotatingCustomLogger(name, facility string, priority syslog.Priority, cConfig, fConfig, sConfig zap.Config, lj *lumberjack.Logger) (*zap.Logger, error) {
	if lj == nil {
		return nil, fmt.Errorf("zaputils: nil lumberjack logger passed")
	}
	return newCustomLogger(name, facility, priority, cConfig, fConfig, sConfig, lj)
}

// ConsoleConfig returns a config that is intended to be used in a console logger.
func ConsoleConfig() zap.Config {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	cfg.ErrorOutputPaths = []string{logFile, "stderr"}
	cfg.OutputPaths = []string{"stdout"}
	return cfg
}

// FileConfig returns a config that is intended to be used in a file logger.
func FileConfig() zap.Config {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.ErrorOutputPaths = []string{logFile, "stderr"}
	cfg.OutputPaths = []string{logFile, "stderr"}
	return cfg
}

// SysConfig returns a config that is inteded to be used in a system logger.
// Note that the Output path is set to an empty slice of strings, because
// this config is meant to be used with a custom writer.
func SysConfig() zap.Config {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.ErrorOutputPaths = []string{logFile, "stderr"}
	cfg.OutputPaths = []string{}
	return cfg
}

func newCustomLogger(name, facility string, priority syslog.Priority, cConfig, fConfig, sConfig zap.Config, lj *lumberjack.Logger) (*zap.Logger, error) {
	cEncoder := zapcore.NewConsoleEncoder(cConfig.EncoderConfig)
	fEncoder := zapcore.NewJSONEncoder(cConfig.EncoderConfig)
	sEncoder := zapcore.NewJSONEncoder(sConfig.EncoderConfig)

	cSink, cCloseOut, err := zap.Open(cConfig.OutputPaths...)
	if err != nil {
		return nil, fmt.Errorf("zaputils: failed to open console config output paths: %w", err)
	}

	defer func() {
		if err != nil {
			cCloseOut()
		}
	}()

	errSink, _, err := zap.Open(cConfig.ErrorOutputPaths...)
	if err != nil {
		return nil, fmt.Errorf("zaputils: failed to open error output paths: %w", err)
	}

	var fSink zapcore.WriteSyncer
	if lj != nil {
		fSink = zapcore.AddSync(lj)
	}

	if fSink == nil {
		var fCloseOut func()
		fSink, fCloseOut, err = zap.Open(fConfig.OutputPaths...)
		if err != nil {
			return nil, fmt.Errorf("zaputils: failed to open file config output paths: %w", err)
		}

		defer func() {
			if err != nil {
				fCloseOut()
			}
		}()
	}

	cCore := zapcore.NewCore(cEncoder, cSink, zap.InfoLevel)
	fCore := zapcore.NewCore(fEncoder, fSink, zap.WarnLevel)
	var sCore zapcore.Core

	if runtime.GOOS != "windows" {
		var sysLogger syslog.Syslogger
		sysLogger, err = syslog.NewLogger(priority, facility, name)
		if err != nil {
			return nil, fmt.Errorf("zaputils: failed to open sys config output paths: %w", err)
		}

		sWriter := zapcore.Lock(zapcore.AddSync(sysLogger))
		sCore = zapcore.NewCore(sEncoder, sWriter, zap.WarnLevel)
	}

	cores := []zapcore.Core{cCore, fCore}
	if sCore != nil {
		cores = append(cores, sCore)
	}

	opts := []zap.Option{zap.ErrorOutput(errSink)}
	core := zapcore.NewTee(cores...)
	return zap.New(core, opts...).Named(name), nil
}
