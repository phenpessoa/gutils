package zaputils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Broker holds 3 zap loggers and provides
// utility methods to work with them.
//
// If created with the NewBroker func:
// P : logs to the console
//
// M : logs to a file
//
// S : logs to systemd - linux only
type Broker struct {
	P *zap.Logger
	M *zap.Logger
	S *zap.Logger
}

// NewBroker creates a new Broker.
// P will be populated with a logger that logs to the console.
// M will be populated with a logger that logs to a file.
// S will be populated wtih a logger that logs to systemd, it will be nil on an windows system.
func NewBroker(name string) (*Broker, error) {
	p, m, s, err := NewCustomLoggers(name)
	if err != nil {
		return nil, err
	}
	return &Broker{P: p, M: m, S: s}, nil
}

// Sync calls the underlying cores' Sync method,
// flushing any buffered log entries.
// Applications should take care to call Sync before exiting.
func (b *Broker) Sync() {
	if b.P != nil {
		_ = b.P.Sync()
	}
	if b.M != nil {
		_ = b.M.Sync()
	}
	if b.S != nil {
		_ = b.S.Sync()
	}
}

// LogToAll logs to all zap's loggers.
// Panic and Fatal levels are not supported.
func (b *Broker) LogToAll(level zapcore.Level, msg string, fields ...zapcore.Field) {
	switch level {
	case zap.DebugLevel, zap.InfoLevel, zap.WarnLevel, zap.ErrorLevel, zap.DPanicLevel:
		if b.P != nil {
			b.P.Log(level, msg, fields...)
		}
		if b.M != nil {
			b.M.Log(level, msg, fields...)
		}
		if b.S != nil {
			b.S.Log(level, msg, fields...)
		}
	}
}
