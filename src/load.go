package main

import (
    "encoding/json"
    log "github.com/sirupsen/logrus"
    "io"
    "io/ioutil"
    "os"
    "os/exec"
    "strings"
    "syscall"
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

type inboundsSettings struct {
    Inbounds []interface{} `json:"inbounds"`
}

type sniffSettings struct {
    Enabled      bool     `json:"enabled"`
    RouteOnly    bool     `json:"routeOnly"`
    DestOverride []string `json:"destOverride"`
}

type inboundSettings struct {
    Tag            string        `json:"tag"`
    Port           int           `json:"port"`
    Protocol       string        `json:"protocol"`
    Settings       interface{}   `json:"settings"`
    StreamSettings interface{}   `json:"streamSettings"`
    Sniffing       sniffSettings `json:"sniffing"`
}

func runCommand(command ...string) (int, string) {
    log.Debugf("Running system command -> %v", command)
    process := exec.Command(command[0], command[1:]...)
    output, _ := process.CombinedOutput()
    log.Debugf("Command %v -> \n%s", command, string(output))
    code := process.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
    if code != 0 {
        log.Warningf("Command %v return code %d", command, code)
    }
    return code, string(output)
}

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
        log.Panicf("Failed to create folder -> %s", folderPath)
    }
}

func listFolder(folderPath string, suffix string) []string {
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

func copyFile(source string, target string) {
    log.Infof("Copy file `%s` => `%s`", source, target)
    srcFile, err := os.Open(source)
    if err != nil {
        log.Panicf("Failed to open file -> %s", source)
    }
    dstFile, err := os.OpenFile(target, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
    if err != nil {
        log.Panicf("Failed to open file -> %s", target)
    }
    _, err = io.Copy(dstFile, srcFile)
    if err != nil {
        log.Panicf("Failed to copy from `%s` to `%s`", source, target)
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
        log.Panicf("File %s -> %v", caption, err)
    }
}

func loadHttpConfig(tag string, port int, sniffObject sniffSettings) interface{} {
    type empty struct{}
    return inboundSettings{
        Tag:            tag,
        Port:           port,
        Protocol:       "http",
        Settings:       empty{},
        StreamSettings: empty{},
        Sniffing:       sniffObject,
    }
}

func loadSocksConfig(tag string, port int, sniffObject sniffSettings) interface{} {
    type empty struct{}
    type socksSettings struct {
        UDP bool `json:"udp"`
    }
    return inboundSettings{
        Tag:            tag,
        Port:           port,
        Protocol:       "socks",
        Settings:       socksSettings{UDP: true},
        StreamSettings: empty{},
        Sniffing:       sniffObject,
    }
}

func loadTProxyConfig(tag string, port int, sniffObject sniffSettings) interface{} {
    type tproxySettings struct {
        Network        string `json:"network"`
        FollowRedirect bool   `json:"followRedirect"`
    }
    type tproxyStreamSettings struct {
        Sockopt struct {
            Tproxy string `json:"tproxy"`
        } `json:"sockopt"`
    }
    tproxyStream := tproxyStreamSettings{}
    tproxyStream.Sockopt.Tproxy = "tproxy"
    return inboundSettings{
        Tag:      tag,
        Port:     port,
        Protocol: "dokodemo-door",
        Settings: tproxySettings{
            Network:        "tcp,udp",
            FollowRedirect: true,
        },
        StreamSettings: tproxyStream,
        Sniffing:       sniffObject,
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

    inbounds := inboundsSettings{}
    sniff := sniffSettings{
        Enabled:      enableSniff,
        RouteOnly:    !enableRedirect,
        DestOverride: []string{"http", "tls"},
    }
    inbounds.Inbounds = append(inbounds.Inbounds, loadTProxyConfig("tproxy", v4TProxyPort, sniff))
    inbounds.Inbounds = append(inbounds.Inbounds, loadTProxyConfig("tproxy6", v6TProxyPort, sniff))
    for tag, port := range httpInbounds {
        inbounds.Inbounds = append(inbounds.Inbounds, loadHttpConfig(tag, port, sniff))
    }
    for tag, port := range socksInbounds {
        inbounds.Inbounds = append(inbounds.Inbounds, loadSocksConfig(tag, port, sniff))
    }
    for _, addon := range addOnInbounds {
        inbounds.Inbounds = append(inbounds.Inbounds, addon)
    }
    inboundsConfig, _ := json.MarshalIndent(inbounds, "", "  ") // json encode
    saveConfig(configDir, "inbounds", string(inboundsConfig)+"\n", true)

    for _, configFile := range listFolder(exposeDir+"/config", ".json") {
        if configFile == "log.json" || configFile == "inbounds.json" {
            log.Warningf("Config file `%s` will be overrided", configFile)
        }
        copyFile(exposeDir+"/config/"+configFile, configDir+"/"+configFile)
    }
}
