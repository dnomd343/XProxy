package logger

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"regexp"
	"strings"
	"testing"
)

var Levels = []Level{DebugLevel, InfoLevel, WarnLevel, ErrorLevel, PanicLevel}

func Test_level(t *testing.T) {
	assert.Equal(t, DebugLevel, zap.DebugLevel)
	assert.Equal(t, InfoLevel, zap.InfoLevel)
	assert.Equal(t, WarnLevel, zap.WarnLevel)
	assert.Equal(t, ErrorLevel, zap.ErrorLevel)
	assert.Equal(t, PanicLevel, zap.PanicLevel)

	for _, level := range Levels {
		SetLevel(level)
		assert.Equal(t, GetLevel(), level)
		assert.Equal(t, logger.verbose, level == DebugLevel) // verbose only for DEBUG level
	}
	SetLevel(InfoLevel) // revert to INFO level
}

func Test_logger(t *testing.T) {
	var plainBuf, plainTeeBuf bytes.Buffer
	var coloredBuf, coloredTeeBuf bytes.Buffer
	logger.plain.writers = []zapcore.WriteSyncer{} // clear slice
	logger.colored.writers = []zapcore.WriteSyncer{}
	AddWriters(false, &plainBuf, &plainTeeBuf)
	AddWriters(true, &coloredBuf, &coloredTeeBuf)
	logger.update() // apply test writers

	printLogs := func(usingLevel Level) {
		SetLevel(usingLevel)
		Debugf("Here is %s level", "DEBUG")
		Infof("Here is %s level", "INFO")
		Warnf("Here is %s level", "WARN")
		Errorf("Here is %s level", "ERROR")
		assert.Panics(t, func() {
			Panicf("Here is %s level", "PANIC")
		})
	}
	printLogs(DebugLevel) // output into buffer
	printLogs(InfoLevel)
	printLogs(WarnLevel)
	printLogs(ErrorLevel)
	printLogs(PanicLevel)

	assertLine := func(log string, level Level, colored bool, verbose bool) {
		regex := `^([\d.:\- ]+) \[(\S+)] (\[\d+] )?\[(\S+)] Here is (\S+) level$`
		if colored {
			regex = `^\x1b\[36m\[XProxy]\x1b\[0m \x1b\[90m([\d.:\- ]+)\x1b\[0m ` +
				`\x1b\[\d\dm\[(\S+)]\x1b\[0m (\x1b\[34m\[\d+]\x1b\[0m )?` +
				`\x1b\[35m\[(\S+)]\x1b\[0m Here is (\S+) level$`
		}
		matches := regexp.MustCompile(regex).FindStringSubmatch(log)
		timeRegex := regexp.MustCompile(`^\d{4}(-\d\d){2} \d{2}(:\d\d){2}\.\d{3}$`)

		assert.NotEmpty(t, matches)                        // valid log line
		assert.Regexp(t, timeRegex, matches[1])            // valid time format
		assert.Equal(t, level.CapitalString(), matches[2]) // valid level string
		assert.Equal(t, level.CapitalString(), matches[5])
		if !verbose {
			assert.Equal(t, "logger/interface_test", matches[4])
		} else {
			assert.Regexp(t, regexp.MustCompile(`^logger/interface_test.go:\d+$`), matches[4])
		}
	}

	assertLogs := func(buffer string, colored bool) {
		var line string
		logs := strings.Split(buffer, "\n")
		for _, limit := range Levels {
			for _, level := range Levels {
				if level >= limit {
					line, logs = logs[0], logs[1:]
					assertLine(line, level, colored, limit == DebugLevel)
				}
			}
		}
		assert.Equal(t, len(logs), 1)
		assert.Equal(t, logs[0], "") // last line is empty
	}

	assertLogs(plainBuf.String(), false)
	assertLogs(coloredBuf.String(), true)
	assert.Equal(t, plainBuf.String(), plainTeeBuf.String())
	assert.Equal(t, coloredBuf.String(), coloredTeeBuf.String())

	logger.plain.writers = []zapcore.WriteSyncer{} // revert logger configure
	logger.colored.writers = []zapcore.WriteSyncer{os.Stdout}
	logger.level.SetLevel(InfoLevel)
	logger.update()
}
