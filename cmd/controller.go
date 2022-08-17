package main

import (
    "XProxy/cmd/asset"
    "XProxy/cmd/config"
    "XProxy/cmd/network"
    "XProxy/cmd/proxy"
)

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

func loadProxy(settings *config.Config) {
    proxy.Load(configDir, exposeDir, proxy.Config{
        Sniff:         settings.EnableSniff,
        Redirect:      settings.EnableRedirect,
        V4TProxyPort:  v4TProxyPort,
        V6TProxyPort:  v6TProxyPort,
        LogLevel:      "warning",
        HttpInbounds:  settings.HttpInbounds,
        SocksInbounds: settings.SocksInbounds,
        AddOnInbounds: settings.AddOnInbounds,
    })
}

func loadAsset(settings *config.Config) {
    asset.LoadGeoSite(assetFile, assetDir)
    asset.LoadGeoIp(assetFile, assetDir)
    asset.AutoUpdate(settings.UpdateCron, settings.UpdateUrls, assetDir)
}
