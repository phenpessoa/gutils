// Package zaputils provides utility funcs to work with
// Uber's zap (https://github.com/uber-go/zap).
//
// A lot of the documentation here is a copy from
// zap's documentation.
//
// Zap is licensed under the MIT LICENSE, and
// so is this package.
//
// Zap's license:
// Copyright (c) 2016-2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package zaputils

import (
	"fmt"
	"runtime"
	"sort"
	"time"

	gsyslog "github.com/hashicorp/go-syslog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	logFile    = "logs.log"
	errLogFile = "logger_errors.log"
)

// NewProdLogger is a helper function that calls NewLogger using NewProdConfig.
//
// See NewProdConfig and NewLogger for details.
func NewProdLogger(name, facility string, priority gsyslog.Priority) (*zap.Logger, func(), error) {
	return NewLogger(NewProdConfig(name, facility, priority))
}

// NewRotatingProdLogger is a helper function that calls NewLogger using NewRotatingProdConfig.
//
// See NewRotatingProdConfig and NewLogger for details.
func NewRotatingProdLogger(name, facility string, priority gsyslog.Priority) (*zap.Logger, func(), error) {
	return NewLogger(NewRotatingProdConfig(name, facility, priority))
}

// NewRotatingProdConfig returns NewProdConfig but with a lumberjackLogger to implement rotation.
//
// NewRotationProdConfig implements a rotating logger for the file logger.
// The maximum size for each file is 100MB, the max age is 14 days,
// the max amount of old log files to keep is 5
// and the old files are compressed with gzip.
//
// See NewProdConfig and NewLogger for details about the rest of the configuration.
func NewRotatingProdConfig(name, facility string, priority gsyslog.Priority) Config {
	cfg := NewProdConfig(name, facility, priority)
	cfg.LumberjackLogger = &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    100,
		MaxAge:     14,
		MaxBackups: 5,
		Compress:   true,
	}
	return cfg
}

// NewProdConfig is a reasonable production logging configuration.
// Console log is enabled at InfoLevel and above, while the file log and the
// sys log are enabled at WarnLevel and above.
//
// It enables sampling and stacktraces are included on logs of ErrorLevel and above.
func NewProdConfig(name, facility string, priority gsyslog.Priority) Config {
	cEncCfg := zap.NewDevelopmentEncoderConfig()
	cEncCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cEncCfg.EncodeTime = zapcore.RFC3339TimeEncoder

	oEncCfg := zap.NewProductionEncoderConfig()
	oEncCfg.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	return Config{
		ConsoleLevel:  zap.NewAtomicLevelAt(zap.InfoLevel),
		FileLevel:     zap.NewAtomicLevelAt(zap.WarnLevel),
		SysLevel:      zap.NewAtomicLevelAt(zap.WarnLevel),
		InitialFields: nil,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		LumberjackLogger:     nil,
		ConsoleEncoderConfig: cEncCfg,
		FileEncoderConfig:    oEncCfg,
		SysEncoderConfig:     oEncCfg,
		SysLogFacility:       facility,
		Name:                 name,
		FileOutputPaths:      []string{logFile},
		ConsoleOutputPaths:   []string{"stderr"},
		ErrorOutputPaths:     []string{"stderr", errLogFile},
		SysLogPriority:       priority,
		Development:          false,
		DisableCaller:        false,
		DisableStacktrace:    false,
	}
}

// Config is used to create a zap Logger.
type Config struct {
	// ConsoleLevel is the Level of the core
	// that logs to the console.
	// If no level is provided, it will
	// default to zap.InfoLevel.
	ConsoleLevel zap.AtomicLevel

	// FileLevel is the Level of the core
	// that logs to the file.
	// If no level is provided, it will
	// default to zap.WarnLevel.
	FileLevel zap.AtomicLevel

	// SysLevel is the Level of the core
	// that logs to the syslog.
	// If no level is provided, it will
	// default to zap.WarnLevel.
	SysLevel zap.AtomicLevel

	// InitialFields is a collection of
	// fields to add to the logger.
	InitialFields map[string]any

	// Sampling sets a sampling policy.
	// A nil SamplingConfig disables sampling.
	Sampling *zap.SamplingConfig

	// LumberjackLogger implements a rotating
	// system to the file logger.
	LumberjackLogger *lumberjack.Logger

	// ConsoleEncoderConfig sets options for
	// the console core encoder.
	// See zapcore.EncoderConfig for details.
	//
	// This will always be console encoded.
	ConsoleEncoderConfig zapcore.EncoderConfig

	// FileEncoderConfig sets options for
	// the file core encoder.
	// See zapcore.EncoderConfig for details.
	//
	// This will always be JSON encoded.
	FileEncoderConfig zapcore.EncoderConfig

	// SysEncoderConfig sets options for
	// the syslog core encoder.
	// See zapcore.EncoderConfig for details.
	//
	// This will always be JSON encoded.
	SysEncoderConfig zapcore.EncoderConfig

	// SysLogFacility represents te syslog
	// facility.
	SysLogFacility string

	// Name is the name of the logger.
	Name string

	// FileOutputPaths is a list of URLs
	// or file paths to write logging
	// output to.
	// It will be used by the file core.
	// See zap.Open for details.
	FileOutputPaths []string

	// FileOutputPaths is a list of URLs
	// or file paths to write logging
	// output to.
	// It will be used by the console core.
	// See zap.Open for details.
	ConsoleOutputPaths []string

	// ErrorOutputPaths is a list of URLs
	// to write internal logger errors to.
	// The default is standard error.
	//
	// Note that this setting only affects
	// internal errors.
	ErrorOutputPaths []string

	// SysLogPriority maps to the syslog
	// priority levels.
	SysLogPriority gsyslog.Priority

	// Development puts the logger in
	// development mode, which changes
	// the behavior of DPanicLevel and
	// takes stacktraces more liberally.
	Development bool

	// DisableCaller stops annotating logs
	// with the calling function's file name
	// and line number. By default, all logs
	// are annotated.
	DisableCaller bool

	// DisableStacktrace completely disables
	// automatic stacktrace capturing.
	// By default, stacktraces are captured
	// for WarnLevel and above logs in development
	// and ErrorLevel and above in production.
	DisableStacktrace bool

	// ExtraCores will be appended to the core tee.
	ExtraCores []zapcore.Core
}

// NewLogger creates a zap Logger that duplicates log entries into 3 or 2 underlying cores.
//
// The func returned by NewLogger calls the underlying Cores's Sync method and closes all opened
// writers. Applications should take care to call this func before exiting.
// This func is never nil and it is safe to be called even if there was an error.
// This func should not be called if the logger is still going to be used.
//
// On windows OS: the logger will have 2 underlying cores. It will log to stderr and to a file
// named "logs.log".
//
// On non windows OS: the logger will have 1 additional core that logs to the syslog.
//
// All logger errors will be logged to a logger_errors.log file and to stderr.
//
// If no level is provided for the ConsoleLevel, it will default to zap.InfoLevel.
// If no level is provided for the FileLevel and/or the SysLevel, they will
// default to zap.WarnLevel.
func NewLogger(c Config) (*zap.Logger, func(), error) {
	outf := func() {}
	if len(c.ConsoleOutputPaths) == 0 {
		return nil, outf, fmt.Errorf("at least one output path for the console core must be specified")
	}

	if len(c.FileOutputPaths) == 0 {
		return nil, outf, fmt.Errorf("at least one output path for the file core must be specified")
	}

	if c.ConsoleLevel == (zap.AtomicLevel{}) {
		c.ConsoleLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	if c.FileLevel == (zap.AtomicLevel{}) {
		c.FileLevel = zap.NewAtomicLevelAt(zap.WarnLevel)
	}

	if c.SysLevel == (zap.AtomicLevel{}) {
		c.SysLevel = zap.NewAtomicLevelAt(zap.WarnLevel)
	}

	if len(c.ErrorOutputPaths) == 0 {
		c.ErrorOutputPaths = []string{"stderr"}
	}

	cEncoder := zapcore.NewConsoleEncoder(c.ConsoleEncoderConfig)
	fEncoder := zapcore.NewJSONEncoder(c.FileEncoderConfig)
	sEncoder := zapcore.NewJSONEncoder(c.SysEncoderConfig)

	cSink, cCloseOut, err := zap.Open(c.ConsoleOutputPaths...)
	if err != nil {
		return nil, outf, fmt.Errorf("zaputils: failed to open console config output paths: %w", err)
	}

	outf = func() { cCloseOut() }

	defer func() {
		if err != nil {
			cCloseOut()
		}
	}()

	errSink, closeErrSink, err := zap.Open(c.ErrorOutputPaths...)
	if err != nil {
		return nil, outf, fmt.Errorf("zaputils: failed to open error output paths: %w", err)
	}

	outf = func() { cCloseOut(); closeErrSink() }

	defer func() {
		if err != nil {
			closeErrSink()
		}
	}()

	var fSink zapcore.WriteSyncer
	if c.LumberjackLogger != nil {
		fSink = zapcore.AddSync(c.LumberjackLogger)
		outf = func() { cCloseOut(); closeErrSink(); c.LumberjackLogger.Close() }
	}

	if fSink == nil {
		var fCloseOut func()
		fSink, fCloseOut, err = zap.Open(c.FileOutputPaths...)
		if err != nil {
			return nil, outf, fmt.Errorf("zaputils: failed to open file config output paths: %w", err)
		}

		outf = func() { cCloseOut(); closeErrSink(); fCloseOut() }

		defer func() {
			if err != nil {
				fCloseOut()
			}
		}()
	}

	cCore := zapcore.NewCore(cEncoder, cSink, c.ConsoleLevel)
	fCore := zapcore.NewCore(fEncoder, fSink, c.FileLevel)
	var sCore zapcore.Core

	if runtime.GOOS != "windows" {
		var sysLogger gsyslog.Syslogger
		sysLogger, err = gsyslog.NewLogger(c.SysLogPriority, c.SysLogFacility, c.Name)
		if err != nil {
			return nil, outf, fmt.Errorf("zaputils: failed to open sys config output paths: %w", err)
		}

		sWriter := zapcore.Lock(zapcore.AddSync(sysLogger))
		sCore = zapcore.NewCore(sEncoder, sWriter, c.SysLevel)
	}

	cores := make([]zapcore.Core, 0, len(c.ExtraCores)+3)
	cores = append(cores, cCore, fCore)
	if sCore != nil {
		cores = append(cores, sCore)
	}
	cores = append(cores, c.ExtraCores...)

	opts := make([]zap.Option, 1, 4)
	opts[0] = zap.ErrorOutput(errSink)

	if !c.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}

	stackLevel := zap.ErrorLevel
	if c.Development {
		stackLevel = zap.WarnLevel
		opts = append(opts, zap.Development())
	}

	if !c.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}

	if scfg := c.Sampling; scfg != nil {
		opts = append(opts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			var samplerOpts []zapcore.SamplerOption
			if scfg.Hook != nil {
				samplerOpts = append(samplerOpts, zapcore.SamplerHook(scfg.Hook))
			}
			return zapcore.NewSamplerWithOptions(
				core,
				time.Second,
				scfg.Initial,
				scfg.Thereafter,
				samplerOpts...,
			)
		}))
	}

	if len(c.InitialFields) > 0 {
		fs := make([]zap.Field, 0, len(c.InitialFields))
		keys := make([]string, 0, len(c.InitialFields))
		for k := range c.InitialFields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fs = append(fs, zap.Any(k, c.InitialFields[k]))
		}
		opts = append(opts, zap.Fields(fs...))
	}

	core := zapcore.NewTee(cores...)
	l := zap.New(core, opts...).Named(c.Name)
	finalf := func() { _ = l.Sync(); outf() }
	return l, finalf, nil
}
