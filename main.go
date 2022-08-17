package main

import (
    "XProxy/cmd/network"
    "XProxy/cmd/proxy"
    "fmt"
    log "github.com/sirupsen/logrus"
)

func main() {
    log.SetLevel(log.DebugLevel)

    fmt.Println("XProxy start")

    network.Load(nil, network.Config{
        RouteTable: 100,
        TProxyPort: 7288,
        Address:    "192.168.2.2",
        Gateway:    "192.168.2.1",
        Bypass:     make([]string, 0),
    }, network.Config{
        RouteTable: 106,
        TProxyPort: 7289,
        Address:    "fc00::2",
        Gateway:    "fc00::1",
        Bypass:     make([]string, 0),
    })

    proxy.Load("/etc/xproxy", "/xproxy", proxy.Config{
        Sniff:         true,
        Redirect:      true,
        V4TProxyPort:  7288,
        V6TProxyPort:  7289,
        LogLevel:      "debug",
        HttpInbounds:  nil,
        SocksInbounds: nil,
        AddOnInbounds: nil,
    })
}
