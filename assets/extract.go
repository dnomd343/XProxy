package assets

import (
	"bytes"
	"compress/bzip2"
	"github.com/gabriel-vasile/mimetype"
	"github.com/klauspost/compress/gzip"
	"github.com/ulikunitz/xz"
	"io"
)

const (
	notArchive = iota
	gzipArchive
	bzip2Archive
	xzArchive
)

// gzipExtract use to extract independent gzip format stream.
func gzipExtract(a *asset) error {
	a.debug("Start extracting gzip archive stream")
	reader, err := gzip.NewReader(a.stream)
	if err != nil {
		a.error("Failed to extract gzip archive -> %v", err)
		return err
	}
	a.stream, a.archive = reader, gzipArchive
	return nil
}

// bzip2Extract use to extract independent bzip2 format stream.
func bzip2Extract(a *asset) error {
	a.debug("Start extracting bzip2 archive stream")
	a.stream = bzip2.NewReader(a.stream)
	a.archive = bzip2Archive
	return nil
}

// xzExtract use to extract independent xz format stream.
func xzExtract(a *asset) error {
	a.debug("Start extracting xz archive stream")
	reader, err := xz.NewReader(a.stream)
	if err != nil {
		a.error("Failed to extract xz archive -> %v", err)
		return err
	}
	a.stream, a.archive = reader, xzArchive
	return nil
}

// tryExtract try to extract the data stream as a compressed format, and will
// return the original data if it cannot be determined.
func (a *asset) tryExtract() error {
	if a.archive != notArchive {
		return nil // already extracted
	}

	header := bytes.NewBuffer(nil)
	mime, err := mimetype.DetectReader(io.TeeReader(a.stream, header))
	if err != nil {
		a.error("Failed to detect data stream -> %v", err)
		return err
	}
	a.stream = io.MultiReader(header, a.stream) // recycle reader

	switch mime.String() { // extract with detected mime type
	case "application/gzip":
		a.debug("Data detected as gzip format")
		return gzipExtract(a)
	case "application/x-bzip2":
		a.debug("Data detected as bzip2 format")
		return bzip2Extract(a)
	case "application/x-xz":
		a.debug("Data detected as xz format")
		return xzExtract(a)
	default:
		a.debug("Data detected as non-archive format -> `%s`", mime)
		return nil
	}
}
