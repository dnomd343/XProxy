package main

import (
    "github.com/robfig/cron"
    log "github.com/sirupsen/logrus"
    "os"
    "os/signal"
    "syscall"
)

var logLevel = "warning"

var v4RouteTable = 100
var v6RouteTable = 106
var v4TProxyPort = 7288
var v6TProxyPort = 7289

var updateCron string
var updateUrls map[string]string

var enableSniff bool
var enableRedirect bool
var httpInbounds map[string]int
var socksInbounds map[string]int
var addOnInbounds []interface{}

var assetFile = "/etc/xproxy/assets.tar.xz"

func main() {
    defer func() {
        if err := recover(); err != nil {
            log.Errorf("Unknown error -> %v", err)
        }
    }()

    log.SetLevel(log.DebugLevel)
    log.Warning("XProxy start")

    content, err := os.ReadFile("test.yml")
    if err != nil {
        panic(err)
    }
    loadConfig(content)

    loadProxy("/etc/xproxy/config", "/xproxy")
    loadGeoSite("/xproxy/assets")
    loadGeoIp("/xproxy/assets")
    autoUpdate := cron.New()
    _ = autoUpdate.AddFunc(updateCron, func() {
        updateAssets("/xproxy/assets")
    })
    autoUpdate.Start()

    loadDns()
    loadNetwork()
    loadTProxy()

    // TODO: running custom script

    xray := newProcess("xray", "-confdir", "/etc/xproxy/config")
    xray.startProcess(true, true)
    subProcess = append(subProcess, xray)
    daemon()

    sigExit := make(chan os.Signal, 1)
    signal.Notify(sigExit, syscall.SIGINT, syscall.SIGTERM)
    <-sigExit
    exit()
}
