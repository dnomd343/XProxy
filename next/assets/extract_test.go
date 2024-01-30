package assets

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"github.com/dsnet/compress/bzip2"
	"github.com/stretchr/testify/assert"
	"github.com/ulikunitz/xz"
	mrand "math/rand"
	"testing"
)

const testMinSize = 16 * 1024 // 16MiB
const testMaxSize = 64 * 1024 // 64MiB

// randBytes generates a specified number of random bytes.
func randBytes(size int) []byte {
	tmp := make([]byte, size)
	_, _ = rand.Read(tmp)
	return tmp
}

func randInt(min int, max int) int {
	return min + int(mrand.Float64()*float64(max-min))
}

func randData() []byte {
	raw := randBytes(1024)
	size := randInt(testMinSize, testMaxSize)
	var buffer bytes.Buffer
	for i := 0; i < size; i++ {
		buffer.Write(raw)
	}
	return buffer.Bytes()
}

func gzipCompress(data []byte) []byte {
	buf := bytes.Buffer{}
	gw := gzip.NewWriter(&buf)
	_, _ = gw.Write(data)
	_ = gw.Close()
	return buf.Bytes()
}

func bzip2Compress(data []byte) []byte {
	buf := bytes.Buffer{}
	bw, _ := bzip2.NewWriter(&buf, &bzip2.WriterConfig{
		Level: bzip2.DefaultCompression,
	})
	_, _ = bw.Write(data)
	_ = bw.Close()
	return buf.Bytes()
}

func xzCompress(data []byte) []byte {
	buf := bytes.Buffer{}
	xw, _ := xz.NewWriter(&buf)
	_, _ = xw.Write(data)
	_ = xw.Close()
	return buf.Bytes()
}

func TestGzipExtract(t *testing.T) {
	raw := randData()
	gzOk := gzipCompress(raw)
	gzErr := append(gzOk, randBytes(randInt(1, 16))...)

	ret, err := gzipExtract(gzOk)
	assert.Nil(t, err)
	assert.Equal(t, raw, ret)
	_, err = gzipExtract(gzErr)
	assert.NotNil(t, err)
}

func TestBzip2Extract(t *testing.T) {
	raw := randData()
	bz2Ok := bzip2Compress(raw)
	bz2Err := append(bz2Ok, randBytes(randInt(1, 16))...)

	ret, err := bzip2Extract(bz2Ok)
	assert.Nil(t, err)
	assert.Equal(t, raw, ret)
	_, err = bzip2Extract(bz2Err)
	assert.NotNil(t, err)
}

func TestXzExtract(t *testing.T) {
	raw := randData()
	xzOk := xzCompress(raw)
	xzErr := append(xzOk, randBytes(randInt(1, 16))...)

	ret, err := xzExtract(xzOk)
	assert.Nil(t, err)
	assert.Equal(t, raw, ret)
	_, err = xzExtract(xzErr)
	assert.NotNil(t, err)
}

func TestExtract(t *testing.T) {
	raw := randData()

	ret, err := tryExtract(raw)
	assert.Nil(t, err)
	assert.Equal(t, raw, ret)

	ret, err = tryExtract(gzipCompress(raw))
	assert.Nil(t, err)
	assert.Equal(t, raw, ret)

	ret, err = tryExtract(bzip2Compress(raw))
	assert.Nil(t, err)
	assert.Equal(t, raw, ret)

	ret, err = tryExtract(xzCompress(raw))
	assert.Nil(t, err)
	assert.Equal(t, raw, ret)
}
