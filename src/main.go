package main

import (
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

    xray := newProcess("xray", "-confdir", "/etc/xproxy/config")
    xray.startProcess(true, true)

    sleep := newProcess("sleep", "1000")
    sleep.startProcess(true, true)

    empty := newProcess("empty")

    subProcess = append(subProcess, xray)
    subProcess = append(subProcess, sleep)
    subProcess = append(subProcess, empty)

    for _, sub := range subProcess {
        daemon(sub)
    }

    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    <-sigs
    exit()

    //content, err := os.ReadFile("test.yml")
    //if err != nil {
    //    panic(err)
    //}
    //loadConfig(content)
    //loadProxy("/etc/xproxy/config", "/xproxy")

    //loadGeoIp("/xproxy/assets")
    //loadGeoSite("/xproxy/assets")
    // TODO: auto-update assets file (by cron command)

    //loadDns()
    //loadNetwork()
    //loadTProxy()

    // TODO: running custom script
    // TODO: start xray service
}
