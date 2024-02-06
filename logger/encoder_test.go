package logger

import (
	"XProxy/mocks"
	"fmt"
	"github.com/magiconair/properties/assert"
	"github.com/petermattis/goid"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zapcore"
	"runtime"
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

func Test_getCaller(t *testing.T) {
	caller := zapCaller()
	assert.Equal(t, getCaller(caller, false), "logger/encoder_test")
	assert.Equal(t, getCaller(caller, true), "logger/encoder_test.go:0")

	caller.File = "Invalid Path"
	assert.Equal(t, getCaller(caller, true), "unknown")
	assert.Equal(t, getCaller(caller, false), "unknown")
}

func Test_timeEncoder(t *testing.T) {
	{
		encoder := mocks.NewPrimitiveArrayEncoder(t)
		encoder.On("AppendString", mock.Anything).Once()
		timeEncoder(testTime, encoder)
		encoder.AssertCalled(t, "AppendString", "2000-01-01 00:00:00.000")
	}
	{
		encoder := mocks.NewPrimitiveArrayEncoder(t)
		encoder.On("AppendString", mock.Anything).Once()
		timeColoredEncoder(testTime, encoder)
		exceptPrefix := "\x1b[36m" + logger.prefix + "\x1b[0m"
		encoder.AssertCalled(t, "AppendString", exceptPrefix+" \x1b[90m2000-01-01 00:00:00.000\x1b[0m")
	}
}

func Test_callerEncoder(t *testing.T) {
	logger.verbose = false
	{
		encoder := mocks.NewPrimitiveArrayEncoder(t)
		encoder.On("AppendString", mock.Anything).Once()
		callerEncoder(zapCaller(), encoder)
		encoder.AssertCalled(t, "AppendString", "[logger/encoder_test]")
	}
	{
		encoder := mocks.NewPrimitiveArrayEncoder(t)
		encoder.On("AppendString", mock.Anything).Once()
		callerColoredEncoder(zapCaller(), encoder)
		encoder.AssertCalled(t, "AppendString", "\x1b[35m[logger/encoder_test]\x1b[0m")
	}

	logger.verbose = true
	{
		encoder := mocks.NewPrimitiveArrayEncoder(t)
		encoder.On("AppendString", mock.Anything).Once()
		callerEncoder(zapCaller(), encoder)
		expectPrefix := fmt.Sprintf("[%d]", goid.Get())
		encoder.AssertCalled(t, "AppendString", expectPrefix+" [logger/encoder_test.go:0]")
	}
	{
		encoder := mocks.NewPrimitiveArrayEncoder(t)
		encoder.On("AppendString", mock.Anything).Once()
		callerColoredEncoder(zapCaller(), encoder)
		expectPrefix := fmt.Sprintf("\x1b[34m[%d]\x1b[0m", goid.Get())
		encoder.AssertCalled(t, "AppendString", expectPrefix+" \x1b[35m[logger/encoder_test.go:0]\x1b[0m")
	}
}

func Test_levelEncoder(t *testing.T) {
	encoder := mocks.NewPrimitiveArrayEncoder(t)
	enroll := func(values []string, call *mock.Call) *mock.Call {
		for _, value := range values {
			if call == nil {
				call = encoder.On("AppendString", value).Once()
			} else {
				call = encoder.On("AppendString", value).Once().NotBefore(call)
			}
		}
		return call
	}

	caller := enroll([]string{
		"[DEBUG]",
		"[INFO]",
		"[WARN]",
		"[ERROR]",
		"[PANIC]",
	}, nil)
	levelEncoder(zapcore.DebugLevel, encoder)
	levelEncoder(zapcore.InfoLevel, encoder)
	levelEncoder(zapcore.WarnLevel, encoder)
	levelEncoder(zapcore.ErrorLevel, encoder)
	levelEncoder(zapcore.PanicLevel, encoder)

	enroll([]string{
		"\x1b[39m[DEBUG]\x1b[0m",
		"\x1b[32m[INFO]\x1b[0m",
		"\x1b[33m[WARN]\x1b[0m",
		"\x1b[31m[ERROR]\x1b[0m",
		"\x1b[91m[PANIC]\x1b[0m",
	}, caller)
	levelColoredEncoder(zapcore.DebugLevel, encoder)
	levelColoredEncoder(zapcore.InfoLevel, encoder)
	levelColoredEncoder(zapcore.WarnLevel, encoder)
	levelColoredEncoder(zapcore.ErrorLevel, encoder)
	levelColoredEncoder(zapcore.PanicLevel, encoder)
}
