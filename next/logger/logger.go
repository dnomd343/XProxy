package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"runtime"
)

type logger struct {
	logger  *zap.Logger
	level   *zap.AtomicLevel
	sugar   *zap.SugaredLogger
	writers []zapcore.WriteSyncer
	stderr  zapcore.Core // fixed stderr output
	prefix  string       // custom output prefix
	path    string       // project absolute path
	verbose bool         // show goroutine id and caller line
}

var logHandle *logger // singleton logger handle

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
	zapLevel := zap.NewAtomicLevelAt(InfoLevel) // using info level in default
	zapCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(logConfig(true)), // colorful output
		zapcore.Lock(os.Stderr),
		zapLevel,
	)
	zapLogger := zap.New(zapCore, zap.AddCaller(), zap.AddCallerSkip(1))
	_, src, _, _ := runtime.Caller(0) // absolute path of current code
	logHandle = &logger{
		logger:  zapLogger,
		level:   &zapLevel,
		stderr:  zapCore,
		sugar:   zapLogger.Sugar(),
		writers: []zapcore.WriteSyncer{},
		path:    path.Join(path.Dir(src), "../"),
		prefix:  "[XProxy]",
		verbose: false,
	}
}

// addWrites adds more plain log writers.
func addWrites(writers ...zapcore.WriteSyncer) {
	logHandle.writers = append(logHandle.writers, writers...)
	plainCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(logConfig(false)), // without colored
		zap.CombineWriteSyncers(logHandle.writers...),
		logHandle.level,
	)
	logHandle.logger = zap.New(
		zapcore.NewTee(logHandle.stderr, plainCore),
		zap.AddCaller(), zap.AddCallerSkip(1),
	)
	logHandle.sugar = logHandle.logger.Sugar()
}
