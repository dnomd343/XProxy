package main

import (
    "fmt"
    log "github.com/sirupsen/logrus"
)

var logLevel = "warning"
var v4TProxyPort = 7288
var v6TProxyPort = 7289

var enableSniff = false
var enableRedirect = true

var httpInbounds = make(map[string]int)
var socksInbounds = make(map[string]int)

func main() {
    log.SetLevel(log.DebugLevel)
    log.Warning("XProxy start")

    httpInbounds["ipv4"] = 1084
    httpInbounds["ipv6"] = 1086
    fmt.Println(httpInbounds)

    socksInbounds["nodeA"] = 1681
    socksInbounds["nodeB"] = 1682
    socksInbounds["nodeC"] = 1683
    fmt.Println(socksInbounds)

    loadProxy("/etc/xproxy/config", "/xproxy")

    //content, err := os.ReadFile("test.yml")
    //if err != nil {
    //	panic(err)
    //}
    //loadConfig(content)
}
