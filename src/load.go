package main

import (
    log "github.com/sirupsen/logrus"
    "os"
    "strings"
)

var logConfig = `{
  "log": {
    "loglevel": "${LEVEL}",
    "access": "${DIR}/access.log",
    "error": "${DIR}/error.log"
  }
}`

var dnsConfig = `{
  "dns": {
    "servers": [
      "localhost"
    ]
  }
}`

var routeConfig = `{
  "routing": {
    "domainStrategy": "AsIs",
    "rules": [
      {
        "type": "field",
        "network": "tcp,udp",
        "outboundTag": "node"
      }
    ]
  }
}`

var outboundsConfig = `{
  "outbounds": [
    {
      "tag": "node",
      "protocol": "freedom",
      "settings": {}
    }
  ]
}`

func isFileExist(filePath string) bool {
    s, err := os.Stat(filePath)
    if err != nil { // file or folder not exist
        return false
    }
    return !s.IsDir()
}

func createFolder(folderPath string) {
    log.Debugf("Loading folder -> %s", folderPath)
    err := os.MkdirAll(folderPath, 0755)
    if err != nil {
        log.Errorf("Create folder `%s` failed", folderPath)
        panic("Create folder failed")
    }
}

func saveConfig(configDir string, caption string, content string, overwrite bool) {
    filePath := configDir + "/" + caption + ".json"
    if !overwrite && isFileExist(filePath) { // file exist and don't overwrite
        log.Debugf("Skip loading config -> %s", filePath)
        return
    }
    log.Debugf("Loading %s -> \n%s", filePath, content)
    err := os.WriteFile(filePath, []byte(content), 0644)
    if err != nil {
        log.Errorf("File %s -> %v", caption, err)
        panic("File save error")
    }
}

func loadProxy(configDir string, exposeDir string) {
    createFolder(exposeDir + "/log")
    createFolder(exposeDir + "/config")
    createFolder(configDir)
    saveConfig(exposeDir+"/config", "dns", dnsConfig+"\n", false)
    saveConfig(exposeDir+"/config", "route", routeConfig+"\n", false)
    saveConfig(exposeDir+"/config", "outbounds", outboundsConfig+"\n", false)

    logConfig = strings.ReplaceAll(logConfig, "${LEVEL}", logLevel)
    logConfig = strings.ReplaceAll(logConfig, "${DIR}", exposeDir+"/log")
    saveConfig(configDir, "log", logConfig+"\n", true)

    // TODO: load inbounds config

    // TODO: copy exposeDir/config/*.json -> configDir (exclude log and inbounds)

}
