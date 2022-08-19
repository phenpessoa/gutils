package zaputils

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestLoggers(t *testing.T) {
	for _, tc := range []struct {
		testName   string
		config     zap.Config
		loggerName string
		shouldErr  bool
	}{
		{
			"empty logger name",
			prettyConfig(),
			"",
			true,
		},
		{
			"invalid config",
			func() zap.Config {
				c := zap.NewDevelopmentConfig()
				c.EncoderConfig.TimeKey = "error"
				c.EncoderConfig.EncodeTime = nil
				return c
			}(),
			"foo",
			true,
		},
		{
			"valid",
			prettyConfig(),
			"foo",
			false,
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			_, err := createLogger(tc.loggerName, tc.config)
			if tc.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateSysLogger(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Log("skipping TestNewCustomLoggers: this test can only be run on unix systems")
		return
	}

	for _, tc := range []struct {
		name       string
		loggerName string
		facility   string
		shouldErr  bool
	}{
		{"empty logger name", "", "USER", true},
		{"invalid facility", "test", "INVALID", true},
		{"valid", "test", "USER", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := createSysLogger(tc.loggerName, tc.facility)
			if tc.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
