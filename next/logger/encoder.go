package logger

import (
	"bytes"
	"fmt"
	"github.com/gookit/color"
	"go.uber.org/zap/zapcore"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// getGid get goroutine ID only for debugging.
// -> https://blog.sgmansfield.com/2015/12/goroutine-ids/
func getGid() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

// getCaller calculate relative source path of caller.
func getCaller(ec zapcore.EntryCaller, verbose bool) string {
	file, err := filepath.Rel(logHandle.path, ec.File)
	if err != nil {
		return "undefined"
	}
	if verbose {
		return file + ":" + strconv.Itoa(ec.Line)
	}
	file, _ = strings.CutSuffix(file, ".go") // remove `.go` suffix
	return file
}

// timeEncoder formats the time as a string.
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// timeColoredEncoder formats the time as a colored string
// with custom prefix.
func timeColoredEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf(
		"%s %s",
		color.Cyan.Render(logHandle.prefix), // colored prefix
		color.Gray.Render(t.Format("2006-01-02 15:04:05.000")),
	))
}

// callerEncoder formats caller in square brackets.
func callerEncoder(ec zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	if !logHandle.verbose {
		enc.AppendString("[" + getCaller(ec, false) + "]")
		return
	}
	enc.AppendString(fmt.Sprintf("[%d] [%s]", getGid(), getCaller(ec, true)))
}

// callerColoredEncoder formats caller in square brackets with
// magenta color.
func callerColoredEncoder(ec zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	if !logHandle.verbose {
		enc.AppendString(color.Magenta.Render("[" + getCaller(ec, false) + "]"))
		return
	}
	enc.AppendString(fmt.Sprintf(
		"%s %s",
		color.Blue.Render(fmt.Sprintf("[%d]", getGid())),
		color.Magenta.Render("["+getCaller(ec, true)+"]"),
	))
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
