package upstream

import (
	"XProxy/logger"
	"io"
	"os"
	"time"
)

type localAsset struct {
	tag  string
	path string
}

func (l *localAsset) Tag() string {
	return l.tag
}

func (l *localAsset) lastModify() time.Time {
	stat, err := os.Stat(l.path)
	if err != nil {
		logger.Warnf("Failed to get local file stat -> %v", err)
		return time.Now() // unknown modify time
	}
	return stat.ModTime() // using last modify time of src file
}

func (l *localAsset) Request() (io.ReadCloser, time.Time, error) {
	logger.Debugf("Start extracting local asset `%s`", l.path)
	stream, err := os.Open(l.path)
	if err != nil {
		logger.Errorf("Failed to read local asset -> %v", err)
		return nil, l.lastModify(), err
	}
	return stream, l.lastModify(), nil
}
