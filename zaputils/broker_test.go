package zaputils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLogToAll(t *testing.T) {
	core, obsp := observer.New(zap.DebugLevel)
	p := zap.New(core, zap.OnFatal(zapcore.WriteThenPanic))

	core, obsm := observer.New(zap.DebugLevel)
	m := zap.New(core, zap.OnFatal(zapcore.WriteThenPanic))

	core, obss := observer.New(zap.DebugLevel)
	s := zap.New(core, zap.OnFatal(zapcore.WriteThenPanic))

	l := &Broker{p, m, s}
	defer l.Sync()

	for _, tc := range []struct {
		name  string
		level zapcore.Level
		msg   string
	}{
		{"debug", zap.DebugLevel, "foo"},
		{"info", zap.InfoLevel, "foo"},
		{"warn", zap.WarnLevel, "foo"},
		{"error", zap.ErrorLevel, "foo"},
		{"dpanic", zap.DPanicLevel, "foo"},
	} {
		// remove previous logs
		obsp.TakeAll()
		obsm.TakeAll()
		obss.TakeAll()

		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)

			l.LogToAll(tc.level, tc.msg)

			assert.Equal(tc.msg, obsp.All()[0].Message)
			assert.Equal(tc.msg, obsm.All()[0].Message)
			assert.Equal(tc.msg, obss.All()[0].Message)
		})
	}
}

func TestNewBroker(t *testing.T) {
	for _, tc := range []struct {
		testName  string
		logName   string
		shouldErr bool
	}{
		{"empty name", "", true},
		{"valid logger", "foo", false},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			_, err := NewBroker(tc.logName)
			if tc.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
