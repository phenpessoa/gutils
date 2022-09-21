package zaputils

import (
	"runtime"
	"testing"

	syslog "github.com/hashicorp/go-syslog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TestLoggers(t *testing.T) {
	assert := assert.New(t)

	_, err := NewProdLogger("prod_logger", "USER", syslog.LOG_DEBUG)
	assert.NoError(err)

	_, err = NewRotatingProdLogger("prod_logger", "USER", syslog.LOG_DEBUG)
	assert.NoError(err)

	_, err = NewRotatingCustomLogger("prod_logger", "USER", syslog.LOG_DEBUG, ConsoleConfig(), FileConfig(), SysConfig(), nil)
	assert.Error(err)

	_, err = NewRotatingCustomLogger("prod_logger", "USER", syslog.LOG_DEBUG, ConsoleConfig(), FileConfig(), SysConfig(), &lumberjack.Logger{})
	assert.NoError(err)
}

func TestCustomLogger(t *testing.T) {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{` ~ ^ ////// \\\\\\\\ invalid path`}

	for _, tc := range []struct {
		name      string
		cConfig   zap.Config
		fConfig   zap.Config
		sConfig   zap.Config
		facility  string
		priority  syslog.Priority
		lj        *lumberjack.Logger
		shouldErr bool
	}{
		{
			"invalid console config",
			func() zap.Config {
				cfg := zap.NewDevelopmentConfig()
				cfg.OutputPaths = []string{` ~ ^ ////// \\\\\\\\ invalid path`}
				return cfg
			}(),
			zap.Config{},
			zap.Config{},
			"",
			syslog.LOG_INFO,
			&lumberjack.Logger{},
			true,
		},
		{
			"invalid error output path",
			func() zap.Config {
				cfg := zap.NewDevelopmentConfig()
				cfg.ErrorOutputPaths = []string{` ~ ^ ////// \\\\\\\\ invalid path`}
				return cfg
			}(),
			zap.Config{},
			zap.Config{},
			"",
			syslog.LOG_INFO,
			&lumberjack.Logger{},
			true,
		},
		{
			"invalid file logger",
			ConsoleConfig(),
			func() zap.Config {
				cfg := zap.NewDevelopmentConfig()
				cfg.OutputPaths = []string{` ~ ^ ////// \\\\\\\\ invalid path`}
				return cfg
			}(),
			zap.Config{},
			"",
			syslog.LOG_INFO,
			nil,
			true,
		},
		{
			"invalid sys logger",
			ConsoleConfig(),
			FileConfig(),
			SysConfig(),
			"INVALID FACILITY",
			syslog.LOG_INFO,
			nil,
			runtime.GOOS != "windows", // this test only tests properly in non windows systems
		},
		{
			"valid",
			ConsoleConfig(),
			FileConfig(),
			SysConfig(),
			"USER",
			syslog.LOG_INFO,
			nil,
			false,
		},
		{
			"valid lj",
			ConsoleConfig(),
			FileConfig(),
			SysConfig(),
			"USER",
			syslog.LOG_INFO,
			&lumberjack.Logger{},
			false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := newCustomLogger("test_logger", tc.facility, tc.priority, tc.cConfig, tc.fConfig, tc.sConfig, tc.lj)
			if tc.shouldErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
