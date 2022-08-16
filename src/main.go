package main

import (
    "fmt"
    log "github.com/sirupsen/logrus"
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
    log.SetLevel(log.DebugLevel)
    log.Warning("XProxy start")

    ls := newProcess("ls", "-al")
    ls.startProcess(true, true)

    fmt.Println(ls.isProcessAlive())
    ls.waitProcess()
    fmt.Println(ls.isProcessAlive())

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
