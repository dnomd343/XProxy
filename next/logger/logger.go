package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Logger struct {
	gid     bool
	logger  *zap.Logger
	level   *zap.AtomicLevel
	sugar   *zap.SugaredLogger
	writers []zapcore.WriteSyncer
	stdCore zapcore.Core
}

var handle Logger

// logConfig generates log config for XProxy.
func logConfig(colored bool) zapcore.EncoderConfig {
	config := zapcore.EncoderConfig{
		ConsoleSeparator: " ",
		MessageKey:       "msg",
		LevelKey:         "level",
		TimeKey:          "time",
		CallerKey:        "caller",
		EncodeTime:       timeEncoder,
		EncodeLevel:      levelEncoder,
		EncodeCaller:     callerEncoder,
	}
	if colored {
		config.EncodeTime = timeColoredEncoder
		config.EncodeLevel = levelColoredEncoder
		config.EncodeCaller = callerColoredEncoder
	}
	return config
}

func init() {
	level := zap.NewAtomicLevelAt(DebugLevel)
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(logConfig(true)),
		zapcore.Lock(os.Stderr), level,
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	handle = Logger{
		gid:     true,
		logger:  logger,
		level:   &level,
		sugar:   logger.Sugar(),
		writers: []zapcore.WriteSyncer{},
		stdCore: core,
	}
}

// addWrites adds more plain log writers.
func addWrites(writers ...zapcore.WriteSyncer) {
	handle.writers = append(handle.writers, writers...)
	plainCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(logConfig(false)),
		zap.CombineWriteSyncers(handle.writers...),
		handle.level,
	)
	handle.logger = zap.New(
		zapcore.NewTee(handle.stdCore, plainCore),
		zap.AddCaller(), zap.AddCallerSkip(1),
	)
	handle.sugar = handle.logger.Sugar()
}
