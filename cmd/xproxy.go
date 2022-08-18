package main

import (
    "XProxy/cmd/config"
    "XProxy/cmd/process"
    "fmt"
    log "github.com/sirupsen/logrus"
    "os"
    "os/signal"
    "syscall"
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

var subProcess []*process.Process

func runProxy() {
    proxy := process.New("xray", "-confdir", configDir)
    proxy.Run(true)
    proxy.Daemon()
    subProcess = append(subProcess, proxy)
}

func blockWait() {
    sigExit := make(chan os.Signal, 1)
    signal.Notify(sigExit, syscall.SIGINT, syscall.SIGTERM) // wait until get exit signal
    <-sigExit
}

func main() {
    defer func() {
        if err := recover(); err != nil {
            log.Errorf("Panic exit -> %v", err)
        }
    }()

    log.SetLevel(log.DebugLevel)
    fmt.Println("XProxy start -> version =", version)

    settings := config.Load(configFile)
    loadNetwork(&settings)
    loadProxy(&settings)
    loadAsset(&settings)
    runScript(&settings)
    runProxy()

    blockWait()
    process.Exit(subProcess...)
}
