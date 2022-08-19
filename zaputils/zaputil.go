package zaputils

import (
	"errors"
	"time"

	syslog "github.com/hashicorp/go-syslog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logFile = "logs.log"
)

var (
	// ErrorEmptyLoggerName will be sent when trying to create a logger with an empty name
	ErrorEmptyLoggerName = errors.New("zaputil: the logger name can not be empty")
)

func prettyConfig() zap.Config {
	prettyConfig := zap.NewDevelopmentConfig()
	prettyConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prettyConfig.DisableStacktrace = true
	prettyConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	prettyConfig.ErrorOutputPaths = []string{logFile}
	prettyConfig.OutputPaths = []string{"stdout"}
	return prettyConfig
}

func mainConfig() zap.Config {
	mainConfig := zap.NewProductionConfig()
	mainConfig.EncoderConfig.TimeKey = "timestamp"
	mainConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	mainConfig.OutputPaths = []string{logFile}
	mainConfig.ErrorOutputPaths = []string{logFile}
	return mainConfig
}

func createLogger(loggerName string, config zap.Config) (*zap.Logger, error) {
	if loggerName == "" {
		return nil, ErrorEmptyLoggerName
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger.Named(loggerName), nil
}

// createSysLogger creates a system logger.
// It will only work on linux systems.
//
// If createSysLogger is called on an windows system,
// the program will fail to build.
func createSysLogger(loggername, facility string) (*zap.Logger, error) {
	if loggername == "" {
		return nil, ErrorEmptyLoggerName
	}

	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	writer, err := syslog.NewLogger(syslog.LOG_INFO, facility, loggername)
	if err != nil {
		return nil, err
	}

	core := zapcore.NewCore(zapcore.NewJSONEncoder(config.EncoderConfig), zapcore.AddSync(writer), config.Level)
	return zap.New(core).Named(loggername), nil
}
