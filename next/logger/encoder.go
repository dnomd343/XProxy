package logger

import (
	"fmt"
	"github.com/gookit/color"
	"go.uber.org/zap/zapcore"
	"time"
)

// timeEncoder formats the time as a string.
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// timeColoredEncoder formats the time as a colored string
// with `[XProxy]` prefix.
func timeColoredEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf(
		"%s %s",
		color.Cyan.Render("[XProxy]"),
		color.Gray.Render(t.Format("2006-01-02 15:04:05.000")),
	))
}

// callerEncoder formats caller in square brackets.
func callerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + caller.TrimmedPath() + "]")
}

// callerColoredEncoder formats caller in square brackets
// with magenta color.
func callerColoredEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(color.Magenta.Render("[" + caller.TrimmedPath() + "]"))
}

// levelEncoder formats log level using square brackets.
func levelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

// levelColoredEncoder formats log level using square brackets
// and uses different colors.
func levelColoredEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	levelStr := "[" + level.CapitalString() + "]"
	switch level {
	case zapcore.DebugLevel:
		levelStr = color.FgDefault.Render(levelStr)
	case zapcore.InfoLevel:
		levelStr = color.Green.Render(levelStr)
	case zapcore.WarnLevel:
		levelStr = color.Yellow.Render(levelStr)
	case zapcore.ErrorLevel:
		levelStr = color.Red.Render(levelStr)
	case zapcore.PanicLevel:
		levelStr = color.LightRed.Render(levelStr)
	}
	enc.AppendString(levelStr)
}
