package config

import (
    "XProxy/cmd/asset"
    "XProxy/cmd/common"
    "XProxy/cmd/custom"
    "XProxy/cmd/dhcp"
    "XProxy/cmd/network"
    "XProxy/cmd/proxy"
    "XProxy/cmd/radvd"
    log "github.com/sirupsen/logrus"
    "os"
    "path"
)

type Config struct {
    Dev    string
    DNS    []string
    IPv4   network.Config
    IPv6   network.Config
    Proxy  proxy.Config
    Asset  asset.Config
    Radvd  radvd.Config
    Custom custom.Config
    DHCP   dhcp.Config
}

func Load(configFile string, config *Config) {
    if !common.IsFileExist(configFile) { // configure not exist -> load default
        loadDefaultConfig(configFile)
    }
    raw, err := os.ReadFile(configFile) // read configure content
    if err != nil {
        log.Panicf("Failed to open %s -> %v", configFile, err)
    }
    rawConfig := configDecode(raw, path.Ext(configFile)) // decode configure content
    decodeDev(&rawConfig, config)
    decodeDns(&rawConfig, config)
    decodeBypass(&rawConfig, config)
    decodeExclude(&rawConfig, config)
    decodeIPv4(&rawConfig, config)
    decodeIPv6(&rawConfig, config)
    decodeProxy(&rawConfig, config)
    decodeUpdate(&rawConfig, config)
    decodeCustom(&rawConfig, config)
    decodeRadvd(&rawConfig, config)
    decodeDhcp(&rawConfig, config)
}
