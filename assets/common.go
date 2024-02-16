package assets

import (
	"XProxy/logger"
	"io"
)

type asset struct {
	tag     string
	stream  io.Reader
	archive int
}

func (a *asset) debug(template string, args ...interface{}) {
	logger.Debugf("[%s] "+template, append([]interface{}{a.tag}, args...)...)
}

func (a *asset) error(template string, args ...interface{}) {
	logger.Errorf("[%s] "+template, append([]interface{}{a.tag}, args...)...)
}

func (a *asset) Read(p []byte) (n int, err error) {
	n, err = a.stream.Read(p)
	if err != nil && err != io.EOF { // data stream broken
		switch a.archive {
		case notArchive:
			a.error("Failed to read data stream -> %v", err)
		case gzipArchive:
			a.error("Failed to extract gzip archive -> %v", err)
		case bzip2Archive:
			a.error("Failed to extract bzip2 archive -> %v", err)
		case xzArchive:
			a.error("Failed to extract xz archive -> %v", err)
		}
	}
	return n, err
}
