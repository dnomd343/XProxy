package logger

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"testing"
)

func Test_init(t *testing.T) {
	pwd := path.Dir(zapCaller().File)
	assert.Equal(t, path.Join(pwd, "../"), project)

	assert.NotNil(t, logger)
	assert.NotNil(t, logger.entry)
	assert.NotNil(t, logger.level)

	assert.Equal(t, logger.verbose, false)
	assert.Equal(t, logger.prefix, "[XProxy]")
	assert.Equal(t, logger.level.Level(), InfoLevel)

	assert.Equal(t, len(logger.plain.writers), 0)
	assert.Equal(t, len(logger.colored.writers), 1)

	for _, level := range []zapcore.Level{DebugLevel, InfoLevel, WarnLevel, ErrorLevel, PanicLevel} {
		logger.level.SetLevel(level)
		assert.Equal(t, logger.entry.Level(), level)
	}
	logger.level.SetLevel(InfoLevel) // revert to INFO level
}

func Test_addWrites(t *testing.T) {
	core := new(logCore)
	assert.Equal(t, len(core.plain.writers), 0)
	assert.Equal(t, len(core.colored.writers), 0)

	core.entry = nil
	core.addPlainWrites(os.Stdout, os.Stdout, os.Stdout)
	assert.NotNil(t, core.entry)
	assert.Equal(t, len(core.plain.writers), 3)

	core.entry = nil
	core.addColoredWrites(os.Stderr, os.Stderr, os.Stderr)
	assert.NotNil(t, core.entry)
	assert.Equal(t, len(core.plain.writers), 3)
}
