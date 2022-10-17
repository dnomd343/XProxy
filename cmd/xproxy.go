package main

import (
    "XProxy/cmd/common"
    "XProxy/cmd/config"
    "XProxy/cmd/custom"
    "XProxy/cmd/process"
    "flag"
    "fmt"
    log "github.com/sirupsen/logrus"
    "io"
    "os"
    "path"
    "runtime"
    "strconv"
)

var version = "1.0.2"
var v4RouteTable = 104
var v6RouteTable = 106
var v4TProxyPort = 7288
var v6TProxyPort = 7289

var proxyBin = ""
var configDir = "/etc/xproxy"
var assetFile = "/assets.tar.xz"

var subProcess []*process.Process
var assetDir, exposeDir, configFile string

func logInit(isDebug bool, logDir string) {
    log.SetFormatter(&log.TextFormatter{
        ForceColors:     true,
        FullTimestamp:   true,
        TimestampFormat: "2006-01-02 15:04:05",
    })
    log.SetLevel(log.InfoLevel) // default log level
    if isDebug {
        log.SetLevel(log.DebugLevel)
    }
    common.CreateFolder(logDir) // confirm log folder exist
    logFile, err := os.OpenFile(path.Join(logDir, "xproxy.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        log.Errorf("Unable to open log file -> %s", path.Join(logDir, "xproxy.log"))
    }
    log.SetOutput(io.MultiWriter(os.Stderr, logFile))
}

func xproxyInit() {
    xproxyConfig := "xproxy.yml"
    if os.Getenv("CONFIG") != "" {
        xproxyConfig = os.Getenv("CONFIG")
    }
    isVersion := flag.Bool("version", false, "Show version")
    isDebug := flag.Bool("debug", os.Getenv("DEBUG") == "true", "Enable debug mode")
    configName := flag.String("config", xproxyConfig, "Config file name")
    flag.Parse()
    if *isVersion { // show version info and exit
        fmt.Printf("XProxy version %s (%s %s/%s)\n", version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
        os.Exit(0)
    }

    exposeDir = "/xproxy" // default folder
    if os.Getenv("EXPOSE_DIR") != "" {
        exposeDir = os.Getenv("EXPOSE_DIR")
    }
    logInit(*isDebug, path.Join(exposeDir, "log"))
    common.CreateFolder(exposeDir)
    assetDir = path.Join(exposeDir, "assets")
    configFile = path.Join(exposeDir, *configName)
    log.Debugf("Expose folder -> %s", exposeDir)
    log.Debugf("Assets folder -> %s", assetDir)
    log.Debugf("Config file -> %s", configFile)

    if os.Getenv("PROXY_BIN") != "" {
        proxyBin = os.Getenv("PROXY_BIN")
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
}

func main() {
    defer func() {
        if err := recover(); err != nil {
            log.Errorf("Panic exit -> %v", err)
        }
    }()
    xproxyInit()
    var settings config.Config
    log.Infof("XProxy %s start (%s %s/%s)", version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
    config.Load(configFile, &settings)
    loadNetwork(&settings)
    loadProxy(&settings)
    loadAsset(&settings)
    loadRadvd(&settings)
    loadDhcp(&settings)
    custom.RunPreScript(&settings.Custom)
    runDhcp(&settings)
    runRadvd(&settings)
    runProxy(&settings)
    blockWait()
    process.Exit(subProcess...)
    custom.RunPostScript(&settings.Custom)
    log.Warningf("All done, goodbye!")
}
