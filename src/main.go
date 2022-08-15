package main

import (
    log "github.com/sirupsen/logrus"
)

var logLevel = "warning"
var v4TProxyPort = 7288
var v6TProxyPort = 7289

var enableSniff = false
var enableRedirect = true

func main() {
    log.SetLevel(log.DebugLevel)
    log.Warning("XProxy start")

    loadProxy("/etc/xproxy/config", "/xproxy")

    //content, err := os.ReadFile("test.yml")
    //if err != nil {
    //	panic(err)
    //}
    //loadConfig(content)
}
