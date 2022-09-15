package zaputils

import (
	"fmt"
	"runtime"
	"time"

	syslog "github.com/hashicorp/go-syslog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
// The system logger only works on non windows systems.
func NewProdLogger(name, facility string, priority syslog.Priority) (*zap.Logger, error) {
	return NewCustomLogger(name, facility, priority, ConsoleConfig(), FileConfig(), SysConfig())
}

// NewCustomLogger creates a new zap logger that consists of 3 sub loggers.
//
// cConfig is the config intended to be used as a console logger.
// The error sink will be grabbed from this config ErrorOutputPaths.
//
// fConfig is the config intended to be used as a file logger. The file path is `logs.log`.
//
// sConfig is the config inteded to be used as a system logger.
// This will only work on non windows systems.
func NewCustomLogger(name, facility string, priority syslog.Priority, cConfig, fConfig, sConfig zap.Config) (*zap.Logger, error) {
	cEncoder := zapcore.NewConsoleEncoder(cConfig.EncoderConfig)
	fEncoder := zapcore.NewJSONEncoder(fConfig.EncoderConfig)
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

	fSink, fCloseOut, err := zap.Open(fConfig.OutputPaths...)
	if err != nil {
		return nil, fmt.Errorf("zaputils: failed to open file config output paths: %w", err)
	}

	defer func() {
		if err != nil {
			fCloseOut()
		}
	}()

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
	cfg.OutputPaths = []string{logFile}
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
