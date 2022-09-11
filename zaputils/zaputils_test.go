package zaputils

import (
	"runtime"
	"testing"

	syslog "github.com/hashicorp/go-syslog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestProdLogger(t *testing.T) {
	_, err := NewProdLogger("prod_logger", "USER", syslog.LOG_DEBUG)
	assert.NoError(t, err)
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
			true,
		},
		{
			"invalid sys logger",
			ConsoleConfig(),
			FileConfig(),
			SysConfig(),
			"INVALID FACILITY",
			syslog.LOG_INFO,
			runtime.GOOS != "windows", // this test only tests properly in non windows systems
		},
		{
			"valid",
			ConsoleConfig(),
			FileConfig(),
			SysConfig(),
			"USER",
			syslog.LOG_INFO,
			false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewCustomLogger("test_logger", tc.facility, tc.priority, tc.cConfig, tc.fConfig, tc.sConfig)
			if tc.shouldErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
