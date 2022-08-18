package config

import (
    "XProxy/cmd/asset"
    "XProxy/cmd/common"
    "XProxy/cmd/network"
    "XProxy/cmd/radvd"
    log "github.com/sirupsen/logrus"
    "os"
)

type Config struct {
    DNS  []string
    IPv4 network.Config
    IPv6 network.Config

    Script   []string
    LogLevel string

    EnableSniff    bool
    EnableRedirect bool
    SniffExclude   []string
    HttpInbounds   map[string]int
    SocksInbounds  map[string]int
    AddOnInbounds  []interface{}

    Update asset.Config
    Radvd  radvd.Config
}

func Load(configFile string) Config {
    if !common.IsFileExist(configFile) { // configure not exist -> load default
        log.Infof("Load default configure -> %s", configFile)
        common.WriteFile(configFile, defaultConfig, false)
    }
    raw, err := os.ReadFile(configFile) // read configure content
    if err != nil {
        log.Panicf("Failed to open %s -> %v", configFile, err)
    }
    rawConfig := yamlDecode(raw) // decode yaml content
    return decode(rawConfig)
}
