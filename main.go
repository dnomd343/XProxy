package main

import (
    "XProxy/cmd/network"
    "fmt"
    log "github.com/sirupsen/logrus"
)

func main() {
    log.SetLevel(log.DebugLevel)

    fmt.Println("xproxy start")
    //common.CreateFolder("/tmp/test")
    //fmt.Println(common.IsFileExist("/tmp/1.jpg"))
    //fmt.Println(common.ListFiles("/xproxy/config", ".json"))
    net = network.Config{
        V4RouteTable: 12,
    }
}
