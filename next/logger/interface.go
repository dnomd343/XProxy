package logger

import (
	"go.uber.org/zap/zapcore"
	"io"
)

const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
	PanicLevel = zapcore.PanicLevel
)

func GetLevel() zapcore.Level {
	return handle.level.Level()
}

func SetLevel(level zapcore.Level) {
	handle.level.SetLevel(level)
}

func Debugf(template string, args ...interface{}) {
	handle.sugar.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	handle.sugar.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	handle.sugar.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	handle.sugar.Errorf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	handle.sugar.Panicf(template, args...)
}

// AddOutputs adds more plain log outputs.
func AddOutputs(outputs ...io.Writer) {
	var writers []zapcore.WriteSyncer
	for _, output := range outputs {
		writers = append(writers, zapcore.AddSync(output))
	}
	addWrites(writers...)
}
