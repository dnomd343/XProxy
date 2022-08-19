package main

import (
    "XProxy/cmd/common"
    "XProxy/cmd/config"
    "XProxy/cmd/process"
    log "github.com/sirupsen/logrus"
    "os"
    "path"
    "strconv"
)

var version = "0.9.2"
var v4RouteTable = 100
var v6RouteTable = 106
var v4TProxyPort = 7288
var v6TProxyPort = 7289
var configDir = "/etc/xproxy"
var assetFile = "/assets.tar.xz"

var goVersion string
var subProcess []*process.Process
var assetDir, exposeDir, configFile string

func xproxyInit() {
    log.SetFormatter(&log.TextFormatter{
        FullTimestamp:   true,
        TimestampFormat: "2006-01-02 15:04:05",
    })
    if len(os.Args) > 1 && os.Args[1] == "--debug" {
        log.SetLevel(log.DebugLevel)
    } else {
        log.SetLevel(log.InfoLevel)
    }

    if os.Getenv("IPV4_TABLE") != "" {
        v4RouteTable, _ = strconv.Atoi(os.Getenv("IPV4_TABLE"))
    }
    if os.Getenv("IPV6_TABLE") != "" {
        v6RouteTable, _ = strconv.Atoi(os.Getenv("IPV6_TABLE"))
    }
    if os.Getenv("IPV4_TPROXY") != "" {
        v4TProxyPort, _ = strconv.Atoi(os.Getenv("IPV4_TPROXY"))
    }
    if os.Getenv("IPV6_TPROXY") != "" {
        v6TProxyPort, _ = strconv.Atoi(os.Getenv("IPV6_TPROXY"))
    }
    log.Debugf("IPv4 Route Table -> %d", v4RouteTable)
    log.Debugf("IPv6 Route Table -> %d", v6RouteTable)
    log.Debugf("IPv4 TProxy Port -> %d", v4TProxyPort)
    log.Debugf("IPv6 TProxy Port -> %d", v6TProxyPort)

    exposeDir = "/xproxy" // default folder
    if os.Getenv("EXPOSE_DIR") != "" {
        exposeDir = os.Getenv("EXPOSE_DIR")
    }
    common.CreateFolder(exposeDir)
    assetDir = path.Join(exposeDir, "assets")
    configFile = path.Join(exposeDir, "xproxy.yml")
    log.Debugf("Expose folder -> %s", exposeDir)
    log.Debugf("Assets folder -> %s", assetDir)
    log.Debugf("Config file -> %s", configFile)
}

func main() {
    defer func() {
        if err := recover(); err != nil {
            log.Errorf("Panic exit -> %v", err)
        }
    }()
    xproxyInit()

    var settings config.Config
    log.Infof("XProxy %s start (%s)", version, goVersion)
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
