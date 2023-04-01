package common

import (
    "github.com/andybalholm/brotli"
    "github.com/go-http-utils/headers"
    "github.com/klauspost/compress/flate"
    "github.com/klauspost/compress/gzip"
    log "github.com/sirupsen/logrus"
    "io"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "strings"
)

func CreateFolder(folderPath string) {
    folder, err := os.Stat(folderPath)
    if err == nil && folder.IsDir() { // folder exist -> skip create process
        return
    }
    log.Debugf("Create folder -> %s", folderPath)
    if err := os.MkdirAll(folderPath, 0755); err != nil {
        log.Errorf("Failed to create folder -> %s", folderPath)
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
    if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
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
        log.Debugf("File `%s` will be overridden", target)
    }
    srcFile, err := os.Open(source)
    defer srcFile.Close()
    if err != nil {
        log.Panicf("Failed to open file -> %s", source)
    }
    dstFile, err := os.OpenFile(target, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
    defer dstFile.Close()
    if err != nil {
        log.Panicf("Failed to open file -> %s", target)
    }
    if _, err = io.Copy(dstFile, srcFile); err != nil {
        log.Panicf("Failed to copy from `%s` to `%s`", source, target)
    }
}

func DownloadBytes(fileUrl string, proxyUrl string) ([]byte, error) {
    client := http.Client{}
    if proxyUrl == "" {
        log.Infof("Download `%s` without proxy", fileUrl)
    } else { // use proxy for download
        log.Infof("Download `%s` via `%s`", fileUrl, proxyUrl)
        rawUrl, _ := url.Parse(proxyUrl)
        client = http.Client{
            Transport: &http.Transport{
                Proxy: http.ProxyURL(rawUrl),
            },
        }
    }

    req, err := http.NewRequest("GET", fileUrl, nil)
    if err != nil {
        log.Errorf("Failed to create http request")
        return nil, err
    }
    req.Header.Set(headers.AcceptEncoding, "gzip, deflate, br")
    resp, err := client.Do(req)
    if err != nil {
        log.Errorf("Failed to execute http GET request")
        return nil, err
    }
    if resp != nil {
        defer resp.Body.Close()
        log.Debugf("Remote data downloaded successfully")
    }

    switch resp.Header.Get(headers.ContentEncoding) {
    case "br":
        log.Debugf("Downloaded content using brolti encoding")
        return io.ReadAll(brotli.NewReader(resp.Body))
    case "gzip":
        log.Debugf("Downloaded content using gzip encoding")
        gr, err := gzip.NewReader(resp.Body)
        if err != nil {
            return nil, err
        }
        return io.ReadAll(gr)
    case "deflate":
        log.Debugf("Downloaded content using deflate encoding")
        zr := flate.NewReader(resp.Body)
        defer zr.Close()
        return io.ReadAll(zr)
    default:
        return io.ReadAll(resp.Body)
    }
}

func DownloadFile(fileUrl string, filePath string, proxyUrl string) bool {
    log.Debugf("File download `%s` => `%s`", fileUrl, filePath)
    data, err := DownloadBytes(fileUrl, proxyUrl)
    if err != nil {
        log.Errorf("Download `%s` error -> %v", fileUrl, err)
        return false
    }

    file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
    defer file.Close()
    if err != nil {
        log.Panicf("Open `%s` error -> %v", filePath, err)
        return false
    }
    if _, err = file.Write(data); err != nil {
        log.Errorf("File `%s` save error -> %v", filePath, err)
        return false
    }
    log.Infof("Download success `%s` => `%s`", fileUrl, filePath)
    return true
}
