package logger

import (
	"github.com/fatih/color"
	"go.uber.org/zap/zapcore"
	"time"
)

func encodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func encodeColoredTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(color.WhiteString(t.Format("2006-01-02 15:04:05.000")))
}

func encodeCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + caller.TrimmedPath() + "]")
}

func encodeColoredCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(color.MagentaString("[" + caller.TrimmedPath() + "]"))
}

func encodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

func encodeColoredLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(func(level zapcore.Level) func(string, ...interface{}) string {
		switch level {
		case zapcore.DebugLevel:
			return color.CyanString
		case zapcore.InfoLevel:
			return color.GreenString
		case zapcore.WarnLevel:
			return color.YellowString
		case zapcore.ErrorLevel:
			return color.RedString
		case zapcore.PanicLevel:
			return color.HiRedString
		default:
			return color.WhiteString
		}
	}(level)("[" + level.CapitalString() + "]"))
}
