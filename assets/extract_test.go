package assets

import (
	"XProxy/logger"
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"fmt"
	"github.com/dsnet/compress/bzip2"
	"github.com/ulikunitz/xz"
	"io"
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
	//size := randInt(testMinSize, testMaxSize)
	size := 257
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

//func TestGzipExtract(t *testing.T) {
//	raw := randData()
//	gzOk := gzipCompress(raw)
//	gzErr := append(gzOk, randBytes(randInt(1, 16))...)
//
//	ret, err := gzipExtract(gzOk)
//	assert.Nil(t, err)
//	assert.Equal(t, raw, ret)
//	_, err = gzipExtract(gzErr)
//	assert.NotNil(t, err)
//}

//func TestBzip2Extract(t *testing.T) {
//	raw := randData()
//	bz2Ok := bzip2Compress(raw)
//	bz2Err := append(bz2Ok, randBytes(randInt(1, 16))...)
//
//	ret, err := bzip2Extract(bz2Ok)
//	assert.Nil(t, err)
//	assert.Equal(t, raw, ret)
//	_, err = bzip2Extract(bz2Err)
//	assert.NotNil(t, err)
//}

//func TestXzExtract(t *testing.T) {
//	raw := randData()
//	xzOk := xzCompress(raw)
//	xzErr := append(xzOk, randBytes(randInt(1, 16))...)
//
//	ret, err := xzExtract(xzOk)
//	assert.Nil(t, err)
//	assert.Equal(t, raw, ret)
//	_, err = xzExtract(xzErr)
//	assert.NotNil(t, err)
//}

//func TestExtract(t *testing.T) {
//	raw := randData()
//
//	ret, err := tryExtract(raw)
//	assert.Nil(t, err)
//	assert.Equal(t, raw, ret)
//
//	ret, err = tryExtract(gzipCompress(raw))
//	assert.Nil(t, err)
//	assert.Equal(t, raw, ret)
//
//	ret, err = tryExtract(bzip2Compress(raw))
//	assert.Nil(t, err)
//	assert.Equal(t, raw, ret)
//
//	ret, err = tryExtract(xzCompress(raw))
//	assert.Nil(t, err)
//	assert.Equal(t, raw, ret)
//}

func Test_demo(t *testing.T) {
	//data := gzipCompress(randData())
	//data := randData()
	//data = append(data, randBytes(randInt(1, 16))...)
	//fmt.Printf("origin gzip size -> %d\n", len(data))

	//data := randData()
	//data := bzip2Compress(randData())
	//data = append(data, randBytes(randInt(1, 16))...)
	//fmt.Printf("origin bzip2 size -> %d\n", len(data))

	//data := randData()
	data := xzCompress(randData())
	data = append(data, randBytes(randInt(1, 16))...)
	fmt.Printf("origin xz size -> %d\n", len(data))

	//buffer := bytes.NewReader(data)

	//reader, err := gzipExtract(buffer)
	//reader, err := bzip2Extract(buffer)
	//reader, err := xzExtract(buffer)

	//archiveType(buffer)

	//fmt.Printf("%v\n", err)
	//
	//buf := make([]byte, 1024*1024*4)
	//for {
	//	n, err := reader.Read(buf)
	//
	//	if err == io.EOF {
	//		fmt.Println("reach stream ending")
	//		break
	//	}
	//	if err != nil {
	//		fmt.Printf("stream error -> %v", err)
	//		return
	//	}
	//
	//	fmt.Printf("get %d bytes\n", n)
	//
	//}

}

func init() {
	logger.SetLevel(logger.DebugLevel)
}

//func Test_archive(t *testing.T) {
//	data := gzipCompress(randData())
//	fmt.Printf("origin gzip size -> %d\n", len(data))
//
//	kk := asset{
//		tag:     "A7B932FD11",
//		stream: bytes.NewReader(data),
//	}
//
//	kk.gzipExtract()
//
//	buf := make([]byte, 4*1024*1024)
//	for {
//		n, err := kk.stream.Read(buf)
//
//		if err == io.EOF {
//			fmt.Printf("get %d bytes\n", n)
//			fmt.Printf("reach stream ending\n")
//			return
//		}
//
//		if err != nil {
//			fmt.Printf("stream error -> %v\n", err)
//			return
//		}
//		fmt.Printf("get %d bytes\n", n)
//	}
//
//}

type brokenReader struct {
	time int
}

func (b *brokenReader) Read(p []byte) (n int, err error) {
	b.time += 1

	fmt.Printf("Read time = %d\n", b.time)

	if b.time < 16 {
		return 1024, nil
	} else {
		return 0, io.ErrShortWrite
	}

}

func Test_extract(t *testing.T) {

	raw := randData()
	data := gzipCompress(raw)
	data = append(data, randBytes(3)...)
	fmt.Printf("origin data size -> %d\n", len(data))

	as := asset{
		tag: "DEMO",
		//stream: &brokenReader{time: 0},
		stream: bytes.NewReader(data),
	}
	if err := as.tryExtract(); err != nil {
		fmt.Printf("try extract error -> %v\n", err)
	} else {
		if n, err := io.Copy(io.Discard, &as); err != nil {
			fmt.Printf("data stream error -> %v\n", err)
		} else {
			fmt.Printf("data stream complete -> %d bytes\n", n)
		}
	}

}
