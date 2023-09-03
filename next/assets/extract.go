package assets

import (
	. "XProxy/next/logger"
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

// gzipExtract use to extract independent gzip archive data.
func gzipExtract(data []byte) ([]byte, error) {
	Logger.Debugf("Start extracting gzip archive -> %d bytes", len(data))
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		Logger.Errorf("Failed to extract gzip archive -> %v", err)
		return nil, err
	}
	defer reader.Close()

	var buffer bytes.Buffer
	size, err := io.Copy(&buffer, reader)
	if err != nil {
		Logger.Errorf("Failed to handle gzip archive -> %v", err)
		return nil, err
	}
	Logger.Debugf("Extracted gzip archive successfully -> %d bytes", size)
	return buffer.Bytes(), nil
}

// bzip2Extract use to extract independent bzip2 archive data.
func bzip2Extract(data []byte) ([]byte, error) {
	Logger.Debugf("Start extracting bzip2 archive -> %d bytes", len(data))
	reader := bzip2.NewReader(bytes.NewReader(data))

	var buffer bytes.Buffer
	size, err := io.Copy(&buffer, reader)
	if err != nil {
		Logger.Errorf("Failed to extract bzip2 archive -> %v", err)
		return nil, err
	}
	Logger.Debugf("Extracted bzip2 archive successfully -> %d bytes", size)
	return buffer.Bytes(), nil
}

// xzExtract use to extract independent xz archive data.
func xzExtract(data []byte) ([]byte, error) {
	Logger.Debugf("Start extracting xz archive -> %d bytes", len(data))
	reader, err := xz.NewReader(bytes.NewReader(data))
	if err != nil {
		Logger.Errorf("Failed to extract xz archive -> %v", err)
		return nil, err
	}

	var buffer bytes.Buffer
	size, err := io.Copy(&buffer, reader)
	if err != nil {
		Logger.Errorf("Failed to handle xz archive -> %v", err)
		return nil, err
	}
	Logger.Debugf("Extracted xz archive successfully -> %d bytes", size)
	return buffer.Bytes(), nil
}

// archiveType use to determine the type of archive file.
func archiveType(data []byte) uint {
	mime := mimetype.Detect(data)
	switch mime.String() {
	case "application/gzip":
		Logger.Debugf("Data detected as gzip format")
		return gzipArchive
	case "application/x-bzip2":
		Logger.Debugf("Data detected as bzip2 format")
		return bzip2Archive
	case "application/x-xz":
		Logger.Debugf("Data detected as xz format")
		return xzArchive
	default:
		Logger.Debugf("Data detected as non-archive format -> `%s`", mime)
		return notArchive
	}
}

// Extract will try to extract the data as a compressed format, and will
// return the original data if it cannot be determined.
func Extract(data []byte) ([]byte, error) {
	switch archiveType(data) {
	case gzipArchive:
		return gzipExtract(data)
	case bzip2Archive:
		return bzip2Extract(data)
	case xzArchive:
		return xzExtract(data)
	default:
		return data, nil
	}
}
