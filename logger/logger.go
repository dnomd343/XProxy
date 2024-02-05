package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"runtime"
)

var project string  // project absolute path
var logger *logCore // singleton logger handle

// logChannel handle multiple writers with unified format.
type logChannel struct {
	encoder zapcore.Encoder
	writers []zapcore.WriteSyncer
}

// logCore manage log level, channels and other interfaces.
type logCore struct {
	prefix  string // custom output prefix
	verbose bool   // show verbose information

	plain   logChannel // log channel with plain text
	colored logChannel // log channel with colored text

	level *zap.AtomicLevel   // zap log level pointer
	entry *zap.SugaredLogger // zap sugared logger entry
}

func init() {
	_, src, _, _ := runtime.Caller(0) // absolute path of current code
	project = path.Join(path.Dir(src), "../")

	zapLevel := zap.NewAtomicLevelAt(InfoLevel) // using info level in default
	logger = &logCore{
		verbose: false,
		level:   &zapLevel,
		prefix:  "[XProxy]",
		plain:   buildChannel(false),
		colored: buildChannel(true),
	}
	logger.addColoredWrites(os.Stderr) // output into stderr in default
}

// buildChannel generate logChannel with `colored` option.
func buildChannel(colored bool) logChannel {
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
	if colored { // using colored version
		config.EncodeTime = timeColoredEncoder
		config.EncodeLevel = levelColoredEncoder
		config.EncodeCaller = callerColoredEncoder
	}
	return logChannel{
		encoder: zapcore.NewConsoleEncoder(config),
		writers: []zapcore.WriteSyncer{}, // without any writer
	}
}

// update refreshes the binding of the log core to the writers.
func (handle *logCore) update() {
	buildCore := func(channel *logChannel) zapcore.Core { // build zap core from logChannel
		return zapcore.NewCore(
			channel.encoder,
			zap.CombineWriteSyncers(channel.writers...),
			handle.level,
		)
	}
	handle.entry = zap.New(
		zapcore.NewTee(buildCore(&handle.plain), buildCore(&handle.colored)),
		zap.AddCaller(), zap.AddCallerSkip(1),
	).Sugar()
}

// addPlainWrites adds plain text writers to the logCore.
func (handle *logCore) addPlainWrites(writers ...zapcore.WriteSyncer) {
	handle.plain.writers = append(handle.plain.writers, writers...)
	handle.update()
}

// addColoredWrites adds colored text writers to the logCore.
func (handle *logCore) addColoredWrites(writers ...zapcore.WriteSyncer) {
	handle.colored.writers = append(handle.colored.writers, writers...)
	handle.update()
}
