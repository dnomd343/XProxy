package main

import (
    "XProxy/cmd/config"
    "XProxy/cmd/process"
    log "github.com/sirupsen/logrus"
)

var version = "0.1.0"

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

func xproxyInit() {
    // log format
    // TODO: set log level
    log.SetLevel(log.DebugLevel)
    // read tproxy port / route table num from env
}

func main() {
    defer func() {
        if err := recover(); err != nil {
            log.Errorf("Panic exit -> %v", err)
        }
    }()
    xproxyInit()

    var settings config.Config
    log.Infof("XProxy %s start", version)
    config.Load(configFile, &settings)
    loadNetwork(&settings)
    loadProxy(&settings)
    loadAsset(&settings)
    loadRadvd(&settings)

    runScript(&settings)
    runProxy(&settings)
    runRadvd(&settings)
    blockWait()
    process.Exit(subProcess...)
}
