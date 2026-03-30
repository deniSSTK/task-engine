package logger

import (
	"libs/env"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

// TODO: add mark secret fields

func NewLogger(defCfg *env.DefConfig) *Logger {
	cfg := zap.NewProductionConfig()

	if defCfg.ENV == env.Dev {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	level, err := zapcore.ParseLevel(env.GetEnv("LOG_LEVEL", "info"))
	if err != nil {
		panic(err)
	}

	cfg.Level = zap.NewAtomicLevelAt(level)
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	base, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return &Logger{Logger: base}
}

func (l *Logger) Named(name string) *Logger {
	return &Logger{Logger: l.Logger.Named(name)}
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{Logger: l.Logger.With(fields...)}
}

func (l *Logger) Sync() error {
	if l == nil || l.Logger == nil {
		return nil
	}

	return l.Logger.Sync()
}
