package logger

import (
	"XProxy/mocks"
	"fmt"
	"github.com/magiconair/properties/assert"
	"github.com/petermattis/goid"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zapcore"
	"runtime"
	"strings"
	"testing"
	"time"
)

var testTime = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

// zapCaller build fake caller with the absolute path of current code.
func zapCaller() zapcore.EntryCaller {
	_, srcPath, _, _ := runtime.Caller(0)
	return zapcore.EntryCaller{
		File: srcPath,
		Line: 0,
	}
}

// encoderTest is a helper function to test buffer output of mock encoder.
func encoderTest(t *testing.T, exec func(*mocks.PrimitiveArrayEncoder), expect string) {
	encoder := mocks.NewPrimitiveArrayEncoder(t)
	encoder.On("AppendInt64", mock.Anything).Maybe() // only enroll used method
	encoder.On("AppendString", mock.Anything).Maybe()

	exec(encoder)
	var values []string
	for _, call := range encoder.Calls {
		values = append(values, fmt.Sprintf("%v", call.Arguments.Get(0)))
	}
	assert.Equal(t, strings.Join(values, " "), expect)
}

func Test_getCaller(t *testing.T) {
	caller := zapCaller()
	caller.File = "Invalid Path"
	assert.Equal(t, getCaller(caller, true), "unknown")
	assert.Equal(t, getCaller(caller, false), "unknown")

	assert.Equal(t, getCaller(zapCaller(), false), "logger/encoder_test")
	assert.Equal(t, getCaller(zapCaller(), true), "logger/encoder_test.go:0")
}

func Test_timeEncoder(t *testing.T) {
	encoderTest(t, func(encoder *mocks.PrimitiveArrayEncoder) {
		timeEncoder(testTime, encoder)
	}, "2000-01-01 00:00:00.000")

	encoderTest(t, func(encoder *mocks.PrimitiveArrayEncoder) {
		timeColoredEncoder(testTime, encoder)
	}, "\x1b[36m"+logger.prefix+"\x1b[0m \x1b[90m2000-01-01 00:00:00.000\x1b[0m")
}

func Test_callerEncoder(t *testing.T) {
	verboseVal := logger.verbose
	callerTest := func(entry func(zapcore.EntryCaller, zapcore.PrimitiveArrayEncoder), expect string) {
		encoderTest(t, func(encoder *mocks.PrimitiveArrayEncoder) {
			entry(zapCaller(), encoder)
		}, expect)
	}

	logger.verbose = false
	callerTest(callerEncoder, "[logger/encoder_test]")
	callerTest(callerColoredEncoder, "\x1b[35m[logger/encoder_test]\x1b[0m")

	logger.verbose = true
	gid := fmt.Sprintf("[%d]", goid.Get())
	callerTest(callerEncoder, gid+" [logger/encoder_test.go:0]")
	callerTest(callerColoredEncoder, "\x1b[34m"+gid+"\x1b[0m \x1b[35m[logger/encoder_test.go:0]\x1b[0m")

	logger.verbose = verboseVal
}

func Test_levelEncoder(t *testing.T) {
	levelTest := func(entry func(zapcore.Level, zapcore.PrimitiveArrayEncoder), level zapcore.Level, expect string) {
		encoderTest(t, func(encoder *mocks.PrimitiveArrayEncoder) {
			entry(level, encoder)
		}, expect)
	}

	levelTest(levelEncoder, zapcore.DebugLevel, "[DEBUG]")
	levelTest(levelEncoder, zapcore.InfoLevel, "[INFO]")
	levelTest(levelEncoder, zapcore.WarnLevel, "[WARN]")
	levelTest(levelEncoder, zapcore.ErrorLevel, "[ERROR]")
	levelTest(levelEncoder, zapcore.PanicLevel, "[PANIC]")

	levelTest(levelColoredEncoder, zapcore.DebugLevel, "\x1b[39m[DEBUG]\x1b[0m")
	levelTest(levelColoredEncoder, zapcore.InfoLevel, "\x1b[32m[INFO]\x1b[0m")
	levelTest(levelColoredEncoder, zapcore.WarnLevel, "\x1b[33m[WARN]\x1b[0m")
	levelTest(levelColoredEncoder, zapcore.ErrorLevel, "\x1b[31m[ERROR]\x1b[0m")
	levelTest(levelColoredEncoder, zapcore.PanicLevel, "\x1b[91m[PANIC]\x1b[0m")
}
