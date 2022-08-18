package proxy

import (
    "XProxy/cmd/common"
)

type Config struct {
    Sniff         bool
    Redirect      bool
    V4TProxyPort  int
    V6TProxyPort  int
    LogLevel      string
    SniffExclude  []string
    HttpInbounds  map[string]int
    SocksInbounds map[string]int
    AddOnInbounds []interface{}
}

func saveConfig(configDir string, caption string, content string, overwrite bool) {
    filePath := configDir + "/" + caption + ".json"
    common.WriteFile(filePath, content+"\n", overwrite)
}

func loadInbounds(config Config) string {
    sniff := sniffObject{
        Enabled:         config.Sniff,
        RouteOnly:       !config.Redirect,
        DestOverride:    []string{"http", "tls", "quic"},
        DomainsExcluded: config.SniffExclude,
    }
    var inbounds []interface{}
    inbounds = append(inbounds, loadTProxyConfig("tproxy", config.V4TProxyPort, sniff))
    inbounds = append(inbounds, loadTProxyConfig("tproxy6", config.V6TProxyPort, sniff))
    for tag, port := range config.HttpInbounds {
        inbounds = append(inbounds, loadHttpConfig(tag, port, sniff))
    }
    for tag, port := range config.SocksInbounds {
        inbounds = append(inbounds, loadSocksConfig(tag, port, sniff))
    }
    for _, addon := range config.AddOnInbounds {
        inbounds = append(inbounds, addon)
    }
    return common.JsonEncode(inboundsObject{
        Inbounds: inbounds,
    })
}

func Load(configDir string, exposeDir string, config Config) {
    common.CreateFolder(exposeDir + "/log")
    common.CreateFolder(exposeDir + "/config")
    common.CreateFolder(configDir)
    saveConfig(exposeDir+"/config", "dns", dnsConfig, false)
    saveConfig(exposeDir+"/config", "route", routeConfig, false)
    saveConfig(exposeDir+"/config", "outbounds", outboundsConfig, false)
    saveConfig(configDir, "inbounds", loadInbounds(config), true)
    saveConfig(configDir, "log", loadLogConfig(config.LogLevel, exposeDir+"/log"), true)
    for _, configFile := range common.ListFiles(exposeDir+"/config", ".json") {
        common.CopyFile(exposeDir+"/config/"+configFile, configDir+"/"+configFile)
    }
}
