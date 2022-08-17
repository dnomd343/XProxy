package common

import (
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func CreateFolder(folderPath string) {
	log.Debugf("Create folder -> %s", folderPath)
	err := os.MkdirAll(folderPath, 0755)
	if err != nil {
		log.Panicf("Failed to create folder -> %s", folderPath)
	}
}

func IsFileExist(filePath string) bool {
	s, err := os.Stat(filePath)
	if err != nil { // file or folder not exist
		return false
	}
	return !s.IsDir()
}

func WriteFile(filePath string, content string, overwrite bool) {
	if !overwrite && IsFileExist(filePath) { // file exist and don't overwrite
		log.Debugf("File `%s` exist -> skip write", filePath)
		return
	}
	log.Debugf("Write file `%s` -> \n%s", filePath, content)
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		log.Panicf("Failed to write `%s` -> %v", filePath, err)
	}
}

func ListFiles(folderPath string, suffix string) []string {
	var fileList []string
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		log.Panicf("Failed to list folder -> %s", folderPath)
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), suffix) {
			fileList = append(fileList, file.Name())
		}
	}
	return fileList
}

func CopyFile(source string, target string) {
	log.Infof("Copy file `%s` => `%s`", source, target)
	if IsFileExist(target) {
		log.Warningf("File `%s` will be overrided", target)
	}
	srcFile, err := os.Open(source)
	if err != nil {
		log.Panicf("Failed to open file -> %s", source)
	}
	dstFile, err := os.OpenFile(target, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		_ = srcFile.Close()
		log.Panicf("Failed to open file -> %s", target)
	}
	_, err = io.Copy(dstFile, srcFile)
	_ = srcFile.Close()
	_ = dstFile.Close()
	if err != nil {
		log.Panicf("Failed to copy from `%s` to `%s`", source, target)
	}
}

func DownloadFile(url string, file string) {
	log.Debugf("File download `%s` => `%s`", url, file)
	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("Download `%s` error -> %v", url, err)
		return
	}
	out, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		_ = resp.Body.Close()
		log.Panicf("Open `%s` error -> %v", file, err)
	}
	_, err = io.Copy(out, resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		log.Panicf("File `%s` save error -> %v", file, err)
	}
	log.Infof("Download success `%s` => `%s`", url, file)
}
