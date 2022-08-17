package main

import (
    "XProxy/cmd/config"
    "fmt"
    log "github.com/sirupsen/logrus"
)

var version = "dev"

var v4RouteTable = 100
var v6RouteTable = 106
var v4TProxyPort = 7288
var v6TProxyPort = 7289

var exposeDir = "/xproxy"
var configDir = "/etc/xproxy"
var assetFile = "/assets.tar.xz"
var assetDir = exposeDir + "/assets"
var configFile = exposeDir + "/config.yml"

func main() {
    log.SetLevel(log.DebugLevel)
    fmt.Println("XProxy start -> version =", version)

    settings := config.Load(configFile)
    loadNetwork(&settings)
    loadProxy(&settings)
    loadAsset(&settings)
    runScript(&settings)

    // TODO: start xray service
}
