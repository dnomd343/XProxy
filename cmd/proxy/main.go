package proxy

import (
    "XProxy/cmd/common"
    log "github.com/sirupsen/logrus"
    "path"
)

type Config struct {
    Log   string         `yaml:"log" json:"json"`
    Core  string         `yaml:"core" json:"core"`
    Http  map[string]int `yaml:"http" json:"http"`
    Socks map[string]int `yaml:"socks" json:"socks"`
    AddOn []interface{}  `yaml:"addon" json:"addon"`
    Sniff struct {
        Enable   bool     `yaml:"enable" json:"enable"`
        Redirect bool     `yaml:"redirect" json:"redirect"`
        Exclude  []string `yaml:"exclude" json:"exclude"`
    } `yaml:"sniff" json:"sniff"`
    V4TProxyPort int
    V6TProxyPort int
}

func saveConfig(configDir string, caption string, content string, overwrite bool) {
    filePath := path.Join(configDir, caption+".json")
    common.WriteFile(filePath, content+"\n", overwrite)
}

func loadInbounds(config *Config) string {
    sniff := sniffObject{
        Enabled:         config.Sniff.Enable,
        RouteOnly:       !config.Sniff.Redirect,
        DestOverride:    []string{"http", "tls", "quic"},
        DomainsExcluded: config.Sniff.Exclude,
    }
    if config.Core == "v2ray" { // PATCH: v2fly-core v4 not support quic sniff
        sniff.DestOverride = sniff.DestOverride[:len(sniff.DestOverride)-1]
    }
    var inbounds []interface{}
    inbounds = append(inbounds, loadTProxyConfig("tproxy", config.V4TProxyPort, sniff))
    inbounds = append(inbounds, loadTProxyConfig("tproxy6", config.V6TProxyPort, sniff))
    for tag, port := range config.Http {
        inbounds = append(inbounds, loadHttpConfig(tag, port, sniff))
    }
    for tag, port := range config.Socks {
        inbounds = append(inbounds, loadSocksConfig(tag, port, sniff))
    }
    for _, addon := range config.AddOn {
        inbounds = append(inbounds, addon)
    }
    return common.JsonEncode(inboundsObject{
        Inbounds: inbounds,
    })
}

func Load(configDir string, exposeDir string, config *Config) {
    common.CreateFolder(path.Join(exposeDir, "log"))
    common.CreateFolder(path.Join(exposeDir, "config"))
    common.CreateFolder(configDir)
    saveConfig(path.Join(exposeDir, "config"), "outbounds", outboundsConfig, false)
    saveConfig(configDir, "inbounds", loadInbounds(config), true)
    saveConfig(configDir, "log", loadLogConfig(config.Log, path.Join(exposeDir, "log")), true)
    for _, configFile := range common.ListFiles(path.Join(exposeDir, "config"), ".json") {
        if configFile == "log.json" || configFile == "inbounds" {
            log.Warningf("Config file %s will be override", configFile)
        }
        common.CopyFile(path.Join(exposeDir, "config", configFile), path.Join(configDir, configFile))
    }
}
