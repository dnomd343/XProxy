package main

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    "syscall"
    "time"
)

var xray *Process
var sleep *Process
var empty *Process

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

func exit() {
    log.Warningf("Start exit process")
    xray.disableProcess()
    xray.sendSignal(syscall.SIGTERM)
    //log.Infof("Send kill signal to process %s", xray.caption)
    sleep.disableProcess()
    sleep.sendSignal(syscall.SIGTERM)
    empty.disableProcess()
    empty.sendSignal(syscall.SIGTERM)
    log.Info("Wait sub process exit")
    for !(xray.done && sleep.done) {
    }
    log.Infof("Exit complete")
}

func main() {
    log.SetLevel(log.DebugLevel)
    log.Warning("XProxy start")

    xray = newProcess("xray", "-confdir", "/etc/xproxy/config")
    xray.startProcess(true, true)

    sleep = newProcess("sleep", "1000")
    sleep.startProcess(true, true)

    //done := make(chan bool, 1)

    daemon(xray)
    daemon(sleep)
    daemon(empty)

    fmt.Println("start sleep...")
    time.Sleep(10 * time.Second)
    fmt.Println("wake up")
    exit()
    
    //<-done

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
