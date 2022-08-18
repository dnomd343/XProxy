package main

import (
    "XProxy/cmd/asset"
    "XProxy/cmd/common"
    "XProxy/cmd/config"
    "XProxy/cmd/network"
    "XProxy/cmd/proxy"
    log "github.com/sirupsen/logrus"
)

func loadAsset(settings *config.Config) {
    asset.Load(assetFile, assetDir)
    asset.AutoUpdate(&settings.Update, assetDir)
}

func loadNetwork(settings *config.Config) {
    settings.IPv4.RouteTable = v4RouteTable
    settings.IPv4.TProxyPort = v4TProxyPort
    settings.IPv6.RouteTable = v6RouteTable
    settings.IPv6.TProxyPort = v6TProxyPort
    network.Load(settings.DNS, settings.IPv4, settings.IPv6)
}

func loadProxy(settings *config.Config) {
    proxy.Load(configDir, exposeDir, proxy.Config{
        Sniff:         settings.EnableSniff,
        Redirect:      settings.EnableRedirect,
        V4TProxyPort:  v4TProxyPort,
        V6TProxyPort:  v6TProxyPort,
        LogLevel:      settings.LogLevel,
        HttpInbounds:  settings.HttpInbounds,
        SocksInbounds: settings.SocksInbounds,
        AddOnInbounds: settings.AddOnInbounds,
    })
}

func runScript(settings *config.Config) {
    for _, script := range settings.Script {
        log.Infof("Run script command -> %s", script)
        common.RunCommand("sh", "-c", script)
    }
}
