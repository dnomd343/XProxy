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

func Debugf(template string, args ...interface{}) {
	logHandle.sugar.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	logHandle.sugar.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	logHandle.sugar.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	logHandle.sugar.Errorf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	logHandle.sugar.Panicf(template, args...)
}

// GetLevel return the current logger level.
func GetLevel() zapcore.Level {
	return logHandle.level.Level()
}

// SetLevel configure logger output level. Note that debug level
// will output more information and reduce performance.
func SetLevel(level zapcore.Level) {
	logHandle.level.SetLevel(level)
	if level == DebugLevel {
		logHandle.verbose = true
	} else {
		logHandle.verbose = false
	}
}

// AddOutputs adds more plain output channel to the logger module.
func AddOutputs(outputs ...io.Writer) {
	var writers []zapcore.WriteSyncer
	for _, output := range outputs {
		writers = append(writers, zapcore.AddSync(output))
	}
	addWrites(writers...)
}
