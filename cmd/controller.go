package main

import (
    "XProxy/cmd/asset"
    "XProxy/cmd/common"
    "XProxy/cmd/config"
    "XProxy/cmd/network"
    "XProxy/cmd/process"
    "XProxy/cmd/proxy"
    "XProxy/cmd/radvd"
    log "github.com/sirupsen/logrus"
    "os"
    "os/signal"
    "path"
    "strconv"
    "syscall"
    "time"
)

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

func loadRadvd(settings *config.Config) {
    radvd.Load(&settings.Radvd)
}

func loadAsset(settings *config.Config) {
    asset.Load(assetFile, assetDir)
    asset.AutoUpdate(&settings.Update, assetDir)
}

func loadNetwork(settings *config.Config) {
    settings.IPv4.RouteTable = v4RouteTable
    settings.IPv4.TProxyPort = v4TProxyPort
    settings.IPv6.RouteTable = v6RouteTable
    settings.IPv6.TProxyPort = v6TProxyPort
    network.Load(settings.DNS, &settings.IPv4, &settings.IPv6)
}

func loadProxy(settings *config.Config) {
    settings.Proxy.V4TProxyPort = v4TProxyPort
    settings.Proxy.V6TProxyPort = v6TProxyPort
    proxy.Load(configDir, exposeDir, &settings.Proxy)
}

func runScript(settings *config.Config) {
    for _, script := range settings.Script {
        log.Infof("Run script command -> %s", script)
        common.RunCommand("sh", "-c", script)
    }
}

func runProxy(settings *config.Config) {
    if settings.Proxy.Core == "xray" { // xray-core
        runProcess("xray", "-confdir", configDir)
    } else if settings.Proxy.Core == "v2ray" { // v2fly-core
        runProcess("v2ray", "-confdir", configDir)
    } else if settings.Proxy.Core == "sagray" { // sager-core
        runProcess("sagray", "run", "-confdir", configDir)
    } else {
        log.Panicf("Unknown core type -> %s", settings.Proxy.Core)
    }
}

func runRadvd(settings *config.Config) {
    if settings.Radvd.Enable {
        radvdCmd := []string{"radvd", "--nodaemon"}
        if settings.Radvd.Log > 0 { // with log option
            radvdCmd = append(radvdCmd, "--logmethod", "logfile")
            radvdCmd = append(radvdCmd, "--logfile", path.Join(exposeDir, "log/radvd.log"))
            radvdCmd = append(radvdCmd, "--debug", strconv.Itoa(settings.Radvd.Log))
            time.Sleep(time.Second) // radvd will crash on first boot without delay (enable debug), why???
        }
        runProcess(radvdCmd...)
    }
}
