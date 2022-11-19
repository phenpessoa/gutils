package zaputils

import (
	"os"
	"runtime"
	"testing"

	gsyslog "github.com/hashicorp/go-syslog"
	"go.uber.org/zap/zapcore"
)

func TestMain(m *testing.M) {
	code := m.Run()
	defer os.Exit(code)
	os.Remove(logFile)
	os.Remove(errLogFile)
}

func TestLoggers(t *testing.T) {
	for _, tc := range []struct {
		f    func() error
		name string
	}{
		{
			f: func() error {
				_, f, err := NewProdLogger("prod_logger", "USER", gsyslog.LOG_DEBUG)
				defer f()
				return err
			},
			name: "NewProdLogger",
		},
		{
			f: func() error {
				_, f, err := NewRotatingProdLogger("prod_logger", "USER", gsyslog.LOG_DEBUG)
				defer f()
				return err
			},
			name: "NewRotatingProdLogger",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.f(); err != nil {
				t.Errorf("%s failed with: %s", tc.name, err)
			}
		})
	}
}

func TestCustomLogger(t *testing.T) {
	for _, tc := range []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name:    "empty console output path",
			cfg:     Config{},
			wantErr: true,
		},
		{
			name: "empty file output path",
			cfg: Config{
				ConsoleOutputPaths: []string{"stderr"},
			},
			wantErr: true,
		},
		{
			name: "invalid console output path",
			cfg: Config{
				ConsoleOutputPaths: []string{` ~ ^ ////// \\\\\\\\ invalid path`},
				FileOutputPaths:    []string{logFile},
			},
			wantErr: true,
		},
		{
			name: "invalid error output path",
			cfg: Config{
				ConsoleOutputPaths: []string{"stderr"},
				FileOutputPaths:    []string{logFile},
				ErrorOutputPaths:   []string{` ~ ^ ////// \\\\\\\\ invalid path`},
			},
			wantErr: true,
		},
		{
			name: "invalid file output path",
			cfg: Config{
				ConsoleOutputPaths: []string{"stderr"},
				FileOutputPaths:    []string{` ~ ^ ////// \\\\\\\\ invalid path`},
			},
			wantErr: true,
		},
		{
			name: "invalid sys log facility",
			cfg: Config{
				ConsoleOutputPaths: []string{"stderr"},
				FileOutputPaths:    []string{logFile},
				SysLogFacility:     "INVALID FACILITY",
			},
			// this test only works on non windows systems
			wantErr: runtime.GOOS != "windows",
		},
		{
			name:    "valid non rotating logger",
			cfg:     NewProdConfig("test_non_rotating_logger", "USER", gsyslog.LOG_INFO),
			wantErr: false,
		},
		{
			name:    "valid rotating logger",
			cfg:     NewRotatingProdConfig("test_rotating_logger", "USER", gsyslog.LOG_INFO),
			wantErr: false,
		},
		{
			name: "valid non rotating pseudo development logger",
			cfg: func() Config {
				cfg := NewProdConfig("test_non_rotating_logger", "USER", gsyslog.LOG_INFO)
				cfg.Development = true
				return cfg
			}(),
			wantErr: false,
		},
		{
			name: "valid non rotating logger with sampling hook",
			cfg: func() Config {
				cfg := NewProdConfig("test_non_rotating_logger", "USER", gsyslog.LOG_INFO)
				cfg.Sampling.Hook = func(e zapcore.Entry, sd zapcore.SamplingDecision) {}
				return cfg
			}(),
			wantErr: false,
		},
		{
			name: "valid non rotating logger with initial fields",
			cfg: func() Config {
				cfg := NewProdConfig("test_non_rotating_logger", "USER", gsyslog.LOG_INFO)
				cfg.InitialFields = map[string]any{"first": 123, "second": "abc"}
				return cfg
			}(),
			wantErr: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, f, err := NewLogger(tc.cfg)
			defer f()
			if (err != nil) != tc.wantErr {
				t.Errorf("NewLogger failed\nWantErr: %v\nErr: %v", tc.wantErr, err)
			}
		})
	}
}
