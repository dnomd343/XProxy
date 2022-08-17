package main

import (
    "XProxy/cmd/config"
    "XProxy/cmd/network"
    "fmt"
    log "github.com/sirupsen/logrus"
)

var exposeDir = "/xproxy"

var v4RouteTable = 100
var v6RouteTable = 106
var v4TProxyPort = 7288
var v6TProxyPort = 7289

func loadNetwork(settings *config.Config) {
    v4Settings := network.Config{
        RouteTable: v4RouteTable,
        TProxyPort: v4TProxyPort,
        Address:    settings.V4Address,
        Gateway:    settings.V4Gateway,
        Bypass:     settings.V4Bypass,
    }
    v6Settings := network.Config{
        RouteTable: v6RouteTable,
        TProxyPort: v6TProxyPort,
        Address:    settings.V6Address,
        Gateway:    settings.V6Gateway,
        Bypass:     settings.V6Bypass,
    }
    network.Load(settings.DNS, v4Settings, v6Settings)
}

func main() {
    log.SetLevel(log.DebugLevel)
    fmt.Println("XProxy start")

    settings := config.Load(exposeDir + "/config.yml")
    fmt.Println(settings)

    loadNetwork(&settings)

    // TODO: load proxy
    // TODO: load asset

    // TODO: start xray service
}
