package main

import (
    "XProxy/cmd/process"
    "fmt"
    log "github.com/sirupsen/logrus"
    "time"
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

    //settings := config.Load(configFile)
    //loadNetwork(&settings)
    //loadProxy(&settings)
    //loadAsset(&settings)
    //runScript(&settings)

    xray := process.New("xray", "-confdir", configDir)
    xray.Run(true)
    xray.Daemon()

    sleep := process.New("sleep", "1001")
    sleep.Run(true)
    sleep.Daemon()

    empty := process.New("empty")
    empty.Daemon()

    time.Sleep(5 * time.Second)

    process.Exit(xray, sleep, empty)

}
