package assets

//const (
//	notArchive = iota
//	gzipArchive
//	bzip2Archive
//	xzArchive
//)

//type archive struct {
//	id    string
//	size  uint64
//	input io.Reader
//}
//
//func (a *archive) Read(p []byte) (n int, err error) {
//	n, err = a.input.Read(p)
//	if err == io.EOF {
//		logger.Debugf("read %d bytes", n)
//
//		logger.Debugf("reach ending")
//		return n, err
//	}
//	if err != nil {
//		logger.Errorf("Failed to extract archive -> %v", err)
//
//		// TODO: do close process
//
//		return n, err
//	}
//
//	logger.Debugf("read %d bytes", n)
//
//	return n, err
//}

// gzipExtract use to extract independent gzip archive stream.
//func gzipExtract(stream io.Reader) (io.Reader, error) {
//	logger.Debugf("Start extracting gzip archive")
//	reader, err := gzip.NewReader(stream)
//	if err != nil {
//		logger.Errorf("Failed to extract gzip archive -> %v", err)
//		return nil, err
//	}
//	//defer reader.Close()
//	return reader, nil
//}

// bzip2Extract use to extract independent bzip2 archive stream.
//func bzip2Extract(stream io.Reader) (io.Reader, error) {
//	logger.Debugf("Start extracting bzip2 archive")
//	reader := bzip2.NewReader(stream)
//	return reader, nil
//}

// xzExtract use to extract independent xz archive stream.
//func xzExtract(stream io.Reader) (io.Reader, error) {
//	logger.Debugf("Start extracting xz archive")
//	reader, err := xz.NewReader(stream)
//	if err != nil {
//		logger.Errorf("Failed to extract xz archive -> %v", err)
//		return nil, err
//	}
//	return reader, nil
//}

// gzipExtract use to extract independent gzip archive data.
//func gzipExtract(data []byte) ([]byte, error) {
//	logger.Debugf("Start extracting gzip archive -> %d bytes", len(data))
//	reader, err := gzip.NewReader(bytes.NewReader(data))
//	if err != nil {
//		logger.Errorf("Failed to extract gzip archive -> %v", err)
//		return nil, err
//	}
//	defer reader.Close()
//
//	var buffer bytes.Buffer
//	size, err := reader.WriteTo(&buffer)
//	if err != nil {
//		logger.Errorf("Failed to handle gzip archive -> %v", err)
//		return nil, err
//	}
//	logger.Debugf("Extracted gzip archive successfully -> %d bytes", size)
//	return buffer.Bytes(), nil
//}

// bzip2Extract use to extract independent bzip2 archive data.
//func bzip2Extract(data []byte) ([]byte, error) {
//	logger.Debugf("Start extracting bzip2 archive -> %d bytes", len(data))
//	reader := bzip2.NewReader(bytes.NewReader(data))
//
//	var buffer bytes.Buffer
//	size, err := io.Copy(&buffer, reader)
//	if err != nil {
//		logger.Errorf("Failed to extract bzip2 archive -> %v", err)
//		return nil, err
//	}
//	logger.Debugf("Extracted bzip2 archive successfully -> %d bytes", size)
//	return buffer.Bytes(), nil
//}

// xzExtract use to extract independent xz archive data.
//func xzExtract(data []byte) ([]byte, error) {
//	logger.Debugf("Start extracting xz archive -> %d bytes", len(data))
//	reader, err := xz.NewReader(bytes.NewReader(data))
//	if err != nil {
//		logger.Errorf("Failed to extract xz archive -> %v", err)
//		return nil, err
//	}
//
//	var buffer bytes.Buffer
//	size, err := io.Copy(&buffer, reader)
//	if err != nil {
//		logger.Errorf("Failed to handle xz archive -> %v", err)
//		return nil, err
//	}
//	logger.Debugf("Extracted xz archive successfully -> %d bytes", size)
//	return buffer.Bytes(), nil
//}

//func recycleReader(input io.Reader) (mimeType string, recycled io.Reader, err error) {
//	// header will store the bytes mimetype uses for detection.
//	header := bytes.NewBuffer(nil)
//
//	// After DetectReader, the data read from input is copied into header.
//	mtype, err := mimetype.DetectReader(io.TeeReader(input, header))
//	if err != nil {
//		fmt.Printf("error -> %v\n", err)
//		return
//	}
//
//	// Concatenate back the header to the rest of the file.
//	// recycled now contains the complete, original data.
//	recycled = io.MultiReader(header, input)
//
//	fmt.Printf("mime-type -> %v\n", mtype)
//
//	return mtype.String(), recycled, err
//}

// archiveType use to determine the type of archive file.
//func archiveType(stream io.Reader) uint {
//
//	fmt.Println("start")
//
//	mime, stream, _ := recycleReader(stream)
//
//	fmt.Println("end")
//
//	//mime := mimetype.Detect(data)
//	switch mime {
//	case "application/gzip":
//		logger.Debugf("Data detected as gzip format")
//		return gzipArchive
//	case "application/x-bzip2":
//		logger.Debugf("Data detected as bzip2 format")
//		return bzip2Archive
//	case "application/x-xz":
//		logger.Debugf("Data detected as xz format")
//		return xzArchive
//	default:
//		logger.Debugf("Data detected as non-archive format -> `%s`", mime)
//		return notArchive
//	}
//}

// archiveType use to determine the type of archive file.
//func archiveType(data []byte) uint {
//	mime := mimetype.Detect(data)
//	switch mime.String() {
//	case "application/gzip":
//		logger.Debugf("Data detected as gzip format")
//		return gzipArchive
//	case "application/x-bzip2":
//		logger.Debugf("Data detected as bzip2 format")
//		return bzip2Archive
//	case "application/x-xz":
//		logger.Debugf("Data detected as xz format")
//		return xzArchive
//	default:
//		logger.Debugf("Data detected as non-archive format -> `%s`", mime)
//		return notArchive
//	}
//}

// tryExtract will try to extract the data as a compressed format, and will
// return the original data if it cannot be determined.
//func tryExtract(data []byte) ([]byte, error) {
//	//switch archiveType(data) {
//	//case gzipArchive:
//	//	//return gzipExtract(data)
//	//case bzip2Archive:
//	//	//return bzip2Extract(data)
//	//case xzArchive:
//	//	//return xzExtract(data)
//	//default:
//	//	return data, nil
//	//}
//	return nil, nil
//}
