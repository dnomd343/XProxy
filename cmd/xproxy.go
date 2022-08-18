package main

import (
    "XProxy/cmd/config"
    "XProxy/cmd/process"
    "XProxy/cmd/radvd"
    log "github.com/sirupsen/logrus"
    "os"
    "os/signal"
    "syscall"
)

var version = "0.0.9"

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

func runProcess(command ...string) {
    sub := process.New(command...)
    sub.Run(true)
    sub.Daemon()
    subProcess = append(subProcess, sub)
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
    log.Infof("XProxy %s start", version)

    settings := config.Load(configFile)
    loadNetwork(&settings)
    loadProxy(&settings)
    loadAsset(&settings)
    radvd.Load(&settings.Radvd)

    runScript(&settings)
    runProcess("xray", "-confdir", configDir)
    if settings.Radvd.Enable {
        runProcess("radvd", "-n", "-m", "logfile", "-l", exposeDir+"/log/radvd.log")
    }
    blockWait()
    process.Exit(subProcess...)
}
