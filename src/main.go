package main

import (
    log "github.com/sirupsen/logrus"
    "os"
)

var logLevel = "warning"
var v4TProxyPort = 7288
var v6TProxyPort = 7289

var enableSniff bool
var enableRedirect bool
var httpInbounds map[string]int
var socksInbounds map[string]int
var addOnInbounds []interface{}

func main() {
    log.SetLevel(log.DebugLevel)
    log.Warning("XProxy start")

    content, err := os.ReadFile("test.yml")
    if err != nil {
        panic(err)
    }
    loadConfig(content)
    loadProxy("/etc/xproxy/config", "/xproxy")
}
