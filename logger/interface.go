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
	logger.entry.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	logger.entry.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	logger.entry.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	logger.entry.Errorf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	logger.entry.Panicf(template, args...)
}

// GetLevel return the current logger level.
func GetLevel() zapcore.Level {
	return logger.level.Level()
}

// SetLevel configure logger output level. Note that debug level will output
// more information and reduce performance.
func SetLevel(level zapcore.Level) {
	logger.level.SetLevel(level)
	if level == DebugLevel {
		logger.verbose = true
	} else {
		logger.verbose = false
	}
}

// AddWriters add more writers to target log channel.
func AddWriters(colored bool, writers ...io.Writer) {
	var syncWriters []zapcore.WriteSyncer
	for _, writer := range writers {
		syncWriters = append(syncWriters, zapcore.AddSync(writer))
	}
	if !colored {
		logger.addPlainWrites(syncWriters...)
	} else {
		logger.addColoredWrites(syncWriters...)
	}
}
