package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
	level  *zap.AtomicLevel
	sugar  *zap.SugaredLogger
}

var handle Logger

func logConfig(level zap.AtomicLevel, colored bool) zapcore.EncoderConfig {
	config := zapcore.EncoderConfig{
		ConsoleSeparator: " ",
		MessageKey:       "msg",
		LevelKey:         "level",
		TimeKey:          "time",
		CallerKey:        "caller",
		EncodeTime:       encodeTime,
		EncodeLevel:      encodeLevel,
		EncodeCaller:     encodeCaller,
	}
	if colored {
		config.EncodeTime = encodeColoredTime
		config.EncodeLevel = encodeColoredLevel
		config.EncodeCaller = encodeColoredCaller
	}
	return config
}

func init() {

	level := zap.NewAtomicLevelAt(DebugLevel)

	writer, _, _ := zap.Open("/dev/stderr", "/root/XProxy/next/lalala.log")

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(logConfig(level, false)),
		//zapcore.Lock(os.Stderr),
		writer,
		level,
	)

	//zapcore.AddSync()
	//zap.Open()

	//return core

	logger := zap.New(core, zap.AddCaller())

	handle = Logger{
		logger: logger,
		level:  &level,
		sugar:  logger.Sugar(),
	}

}
