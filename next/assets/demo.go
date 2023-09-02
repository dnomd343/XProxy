package assets

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

func gzipExtract(reader io.Reader) ([]byte, error) {
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		fmt.Println("gzip content error")
	}
	//
	defer gzipReader.Close()

	//var buffer bytes.Buffer
	//_, err = io.Copy()
	var buffer []byte
	buffer, err = io.ReadAll(gzipReader)
	fmt.Println(len(buffer))

	//fmt.Println(buffer)
	//fmt.Println(string(buffer))

	//gzipFile, err := os.Open("")
	//if err != nil {
	//	return nil, nil
	//}
	//defer gzipFile.Close()
	//gzipReader, err := gzip.NewReader(gzipFile)
	//if err != nil {
	//	return nil, nil
	//}
	//defer gzipReader.Close()
	//var buf bytes.Buffer
	//_, err = io.Copy(&buf, gzipReader)
	//if err != nil {
	//	return nil, err
	//}
	//return buf.Bytes(), nil

	return nil, nil
}

func Demo() {
	fmt.Println("assets demo")

	path := "/root/XProxy/LICENSE.gz"

	fp, err := os.Open(path)
	if err != nil {
		fmt.Println("open failed")
	}
	defer fp.Close()

	//gzipDemo(fp)

	//fmt.Printf("name -> %s\n", fp.Name())

	//var buffer []byte
	//n, err := fp.Read(buffer)
	//fmt.Println(n)
	//fmt.Println(err)
	//buffer, err := io.ReadAll(fp)

	//buffer, err := os.ReadFile(path)
	//if err != nil {
	//	fmt.Printf("error -> %s\n", err)
	//}
	//fmt.Printf("buffer size -> %d\n", len(buffer))
	//fmt.Printf("buffer -> %b\n", buffer)

	//gzipExtract(buffer)
	gzipExtract(fp)

	//fp.Name()

}
