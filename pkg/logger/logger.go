// Package logger provides interface to zap logger
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger extends zap.Logger
type Logger struct {
	*zap.Logger
}

// Config Logger config
type Config struct {
	Debug   bool
	Service string
	Env     string
}

// New creates new logger instance
func New(cfg Config) *Logger {
	var zcfg zap.Config

	if cfg.Debug {
		zcfg = zap.NewDevelopmentConfig()
		zcfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		zcfg = zap.NewProductionConfig()
		zcfg.Encoding = "json"
		zcfg.EncoderConfig.TimeKey = "ts"
		zcfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zcfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	zcfg.OutputPaths = []string{"stderr"}
	zcfg.ErrorOutputPaths = []string{"stderr"}

	z, err := zcfg.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		panic(err)
	}

	fields := make([]zap.Field, 0, 2)
	if cfg.Service != "" {
		fields = append(fields, zap.String("service", cfg.Service))
	}
	if cfg.Env != "" {
		fields = append(fields, zap.String("env", cfg.Env))
	}
	if len(fields) > 0 {
		z = z.With(fields...)
	}

	return &Logger{Logger: z}
}

// Sugar Sugar wraps the Logger to provide a more ergonomic, but slightly slower, API. Sugaring a
// Logger is quite inexpensive, so it's reasonable for a single application to use both Loggers and
// SugaredLoggers, converting between them on the boundaries of performance-sensitive code.
func (l *Logger) Sugar() *zap.SugaredLogger {
	return l.Logger.Sugar()
}

// Named creates new named logger instance
func (l *Logger) Named(name string) *Logger {
	return &Logger{Logger: l.Logger.Named(name)}
}

// With adds defautl fields to logger
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{Logger: l.Logger.With(fields...)}
}
