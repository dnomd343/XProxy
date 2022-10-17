package config

import (
    "XProxy/cmd/asset"
    "XProxy/cmd/common"
    "XProxy/cmd/custom"
    "XProxy/cmd/dhcp"
    "XProxy/cmd/proxy"
    "XProxy/cmd/radvd"
    "encoding/json"
    "github.com/BurntSushi/toml"
    log "github.com/sirupsen/logrus"
    "gopkg.in/yaml.v3"
    "net/url"
)

type NetConfig struct {
    Gateway string `yaml:"gateway" json:"gateway" toml:"gateway"` // network gateway
    Address string `yaml:"address" json:"address" toml:"address"` // network address
}

type RawConfig struct {
    Asset   asset.Config  `yaml:"asset" json:"asset" toml:"asset"`
    Radvd   radvd.Config  `yaml:"radvd" json:"radvd" toml:"radvd"`
    DHCP    dhcp.Config   `yaml:"dhcp" json:"dhcp" toml:"dhcp"`
    Proxy   proxy.Config  `yaml:"proxy" json:"proxy" toml:"proxy"`
    Custom  custom.Config `yaml:"custom" json:"custom" toml:"custom"`
    Network struct {
        Dev     string    `yaml:"dev" json:"dev" toml:"dev"`
        DNS     []string  `yaml:"dns" json:"dns" toml:"dns"`
        ByPass  []string  `yaml:"bypass" json:"bypass" toml:"bypass"`
        Exclude []string  `yaml:"exclude" json:"exclude" toml:"exclude"`
        IPv4    NetConfig `yaml:"ipv4" json:"ipv4" toml:"ipv4"`
        IPv6    NetConfig `yaml:"ipv6" json:"ipv6" toml:"ipv6"`
    } `yaml:"network" json:"network" toml:"network"`
}

func configDecode(raw []byte, fileSuffix string) RawConfig {
    var config RawConfig
    log.Debugf("Config raw content -> \n%s", string(raw))
    if fileSuffix == ".json" {
        if err := json.Unmarshal(raw, &config); err != nil { // json format decode
            log.Panicf("Decode JSON config file error -> %v", err)
        }
    } else if fileSuffix == ".toml" {
        if err := toml.Unmarshal(raw, &config); err != nil { // toml format decode
            log.Panicf("Decode TOML config file error -> %v", err)
        }
    } else {
        if err := yaml.Unmarshal(raw, &config); err != nil { // yaml format decode
            log.Panicf("Decode YAML config file error -> %v", err)
        }
    }
    log.Debugf("Decoded configure -> %v", config)
    return config
}

func decodeDev(rawConfig *RawConfig, config *Config) {
    config.Dev = rawConfig.Network.Dev
    if config.Dev == "" {
        setV4 := rawConfig.Network.IPv4.Address != "" || rawConfig.Network.IPv4.Gateway != ""
        setV6 := rawConfig.Network.IPv6.Address != "" || rawConfig.Network.IPv6.Gateway != ""
        if setV4 || setV6 {
            log.Panicf("Missing dev option in network settings")
        }
    }
    log.Debugf("Network device -> %s", config.Dev)
}

func decodeDns(rawConfig *RawConfig, config *Config) {
    for _, address := range rawConfig.Network.DNS { // dns options
        if common.IsIPv4(address, false) || common.IsIPv6(address, false) {
            config.DNS = append(config.DNS, address)
        } else {
            log.Panicf("Invalid DNS server -> %s", address)
        }
    }
    log.Debugf("DNS server -> %v", config.DNS)
}

func decodeBypass(rawConfig *RawConfig, config *Config) {
    for _, address := range rawConfig.Network.ByPass { // bypass options
        if common.IsIPv4(address, true) || common.IsIPv4(address, false) {
            config.IPv4.Bypass = append(config.IPv4.Bypass, address)
        } else if common.IsIPv6(address, true) || common.IsIPv6(address, false) {
            config.IPv6.Bypass = append(config.IPv6.Bypass, address)
        } else {
            log.Panicf("Invalid bypass IP or CIDR -> %s", address)
        }
    }
    log.Debugf("IPv4 bypass -> %s", config.IPv4.Bypass)
    log.Debugf("IPv6 bypass -> %s", config.IPv6.Bypass)
}

func decodeExclude(rawConfig *RawConfig, config *Config) {
    for _, address := range rawConfig.Network.Exclude { // exclude options
        if common.IsIPv4(address, true) || common.IsIPv4(address, false) {
            config.IPv4.Exclude = append(config.IPv4.Exclude, address)
        } else if common.IsIPv6(address, true) || common.IsIPv6(address, false) {
            config.IPv6.Exclude = append(config.IPv6.Exclude, address)
        } else {
            log.Panicf("Invalid exclude IP or CIDR -> %s", address)
        }
    }
    log.Debugf("IPv4 exclude -> %s", config.IPv4.Exclude)
    log.Debugf("IPv6 exclude -> %s", config.IPv6.Exclude)
}

func decodeIPv4(rawConfig *RawConfig, config *Config) {
    config.IPv4.Address = rawConfig.Network.IPv4.Address
    config.IPv4.Gateway = rawConfig.Network.IPv4.Gateway
    if config.IPv4.Address != "" && !common.IsIPv4(config.IPv4.Address, true) {
        log.Panicf("Invalid IPv4 address (CIDR) -> %s", config.IPv4.Address)
    }
    if config.IPv4.Gateway != "" && !common.IsIPv4(config.IPv4.Gateway, false) {
        log.Panicf("Invalid IPv4 gateway -> %s", config.IPv4.Gateway)
    }
    log.Debugf("IPv4 -> address = %s | gateway = %s", config.IPv4.Address, config.IPv4.Gateway)
}

func decodeIPv6(rawConfig *RawConfig, config *Config) {
    config.IPv6.Address = rawConfig.Network.IPv6.Address
    config.IPv6.Gateway = rawConfig.Network.IPv6.Gateway
    if config.IPv6.Address != "" && !common.IsIPv6(config.IPv6.Address, true) {
        log.Panicf("Invalid IPv6 address (CIDR) -> %s", config.IPv6.Address)
    }
    if config.IPv6.Gateway != "" && !common.IsIPv6(config.IPv6.Gateway, false) {
        log.Panicf("Invalid IPv6 gateway -> %s", config.IPv6.Gateway)
    }
    log.Debugf("IPv6 -> address = %s | gateway = %s", config.IPv6.Address, config.IPv6.Gateway)
}

func decodeProxy(rawConfig *RawConfig, config *Config) {
    config.Proxy = rawConfig.Proxy
    if config.Proxy.Bin == "" {
        config.Proxy.Bin = "xray" // default proxy bin
    }
    log.Debugf("Proxy bin -> %s", config.Proxy.Bin)
    log.Debugf("Proxy log level -> %s", config.Proxy.Log)
    log.Debugf("Http inbounds -> %v", config.Proxy.Http)
    log.Debugf("Socks5 inbounds -> %v", config.Proxy.Socks)
    log.Debugf("Add-on inbounds -> %v", config.Proxy.AddOn)
    log.Debugf("Connection sniff -> %t", config.Proxy.Sniff.Enable)
    log.Debugf("Connection redirect -> %t", config.Proxy.Sniff.Redirect)
    log.Debugf("Connection sniff exclude -> %v", config.Proxy.Sniff.Exclude)
}

func decodeRadvd(rawConfig *RawConfig, config *Config) {
    config.Radvd = rawConfig.Radvd
    if config.Radvd.Enable && config.Radvd.Dev == "" {
        log.Panicf("Radvd enabled without dev option")
    }
    log.Debugf("Radvd log level -> %d", config.Radvd.Log)
    log.Debugf("Radvd network dev -> %s", config.Radvd.Dev)
    log.Debugf("Radvd enable -> %t", config.Radvd.Enable)
    log.Debugf("Radvd options -> %v", config.Radvd.Option)
    log.Debugf("Radvd prefix -> %v", config.Radvd.Prefix)
    log.Debugf("Radvd route -> %v", config.Radvd.Route)
    log.Debugf("Radvd clients -> %v", config.Radvd.Client)
    log.Debugf("Radvd RDNSS -> %v", config.Radvd.RDNSS)
    log.Debugf("Radvd DNSSL -> %v", config.Radvd.DNSSL)
}

func decodeDhcp(rawConfig *RawConfig, config *Config) {
    config.DHCP = rawConfig.DHCP
    log.Debugf("DHCPv4 enable -> %t", config.DHCP.IPv4.Enable)
    log.Debugf("DHCPv4 config -> \n%s", config.DHCP.IPv4.Configure)
    log.Debugf("DHCPv6 enable -> %t", config.DHCP.IPv6.Enable)
    log.Debugf("DHCPv6 config -> \n%s", config.DHCP.IPv6.Configure)
}

func decodeUpdate(rawConfig *RawConfig, config *Config) {
    config.Asset = rawConfig.Asset
    if config.Asset.Update.Proxy != "" {
        _, err := url.Parse(config.Asset.Update.Proxy) // check proxy info
        if err != nil {
            log.Panicf("Invalid asset update proxy -> %s", config.Asset.Update.Proxy)
        }
    }
    log.Debugf("Asset disable -> %t", config.Asset.Disable)
    log.Debugf("Asset update proxy -> %s", config.Asset.Update.Proxy)
    log.Debugf("Asset update cron -> %s", config.Asset.Update.Cron)
    log.Debugf("Asset update urls -> %v", config.Asset.Update.Url)
}

func decodeCustom(rawConfig *RawConfig, config *Config) {
    config.Custom = rawConfig.Custom
    log.Debugf("Custom pre-script -> %v", config.Custom.Pre)
    log.Debugf("Custom post-script -> %v", config.Custom.Post)
}
