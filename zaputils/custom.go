package zaputils

import (
	"runtime"

	"go.uber.org/zap"
)

// NewCustomLoggers creates three custom loggers for the plugin specified.
//
// The sys logger will only be created on non windows systems.
func NewCustomLoggers(pluginName string) (pretty *zap.Logger, main *zap.Logger, sys *zap.Logger, err error) {
	return newCustomLoggers(pluginName, pluginName, pluginName)
}

// newCustomLoggers is the actual worker for NewCustomLoggers.
//
// This is necessary to test.
func newCustomLoggers(prettyName, mainName, sysName string) (pretty *zap.Logger, main *zap.Logger, sys *zap.Logger, err error) {
	pretty, err = createLogger(prettyName, prettyConfig())
	if err != nil {
		return nil, nil, nil, err
	}

	main, err = createLogger(mainName, mainConfig())
	if err != nil {
		return nil, nil, nil, err
	}

	if runtime.GOOS != "windows" {
		sys, err = createSysLogger(sysName, "USER")
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return pretty, main, sys, nil
}
