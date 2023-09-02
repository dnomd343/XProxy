package assets

import (
	. "XProxy/next/logger"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"github.com/ulikunitz/xz"
	"io"
	"os"
)

const gzSample = "/root/XProxy/LICENSE.gz"
const xzSample = "/root/XProxy/LICENSE.xz"
const bz2Sample = "/root/XProxy/LICENSE.bz2"

func gzipExtract(content io.Reader) ([]byte, error) {
	Logger.Debugf("Start extracting gzip archive")
	reader, err := gzip.NewReader(content)
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
	Logger.Debugf("Successfully extracted gzip archive -> %d bytes", size)
	return buffer.Bytes(), nil
}

func bzip2Extract(content io.Reader) ([]byte, error) {
	Logger.Debugf("Start extracting bzip2 archive")
	reader := bzip2.NewReader(content)

	var buffer bytes.Buffer
	size, err := io.Copy(&buffer, reader)
	if err != nil {
		Logger.Errorf("Failed to extract bzip2 archive -> %v", err)
		return nil, err
	}
	Logger.Debugf("Successfully extracted bzip2 archive -> %d bytes", size)
	return buffer.Bytes(), nil
}

func xzExtract(content io.Reader) ([]byte, error) {
	Logger.Debugf("Start extracting xz archive")
	reader, err := xz.NewReader(content)
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
	Logger.Debugf("Successfully extracted xz archive -> %d bytes", size)
	return buffer.Bytes(), nil
}

func Demo() {
	Logger.Infof("Assets demo begin")

	//fp, err := os.Open(gzSample)
	//fp, err := os.Open(bz2Sample)
	fp, err := os.Open(xzSample)
	if err != nil {
		fmt.Println("open failed")
	}
	defer fp.Close()

	//gzipExtract(fp)
	//bzip2Extract(fp)
	xzExtract(fp)

	//fp.Name()

}
