package config

import (
    "XProxy/cmd/asset"
    "XProxy/cmd/common"
    "XProxy/cmd/proxy"
    "XProxy/cmd/radvd"
    log "github.com/sirupsen/logrus"
    "gopkg.in/yaml.v3"
)

type yamlNetConfig struct {
    Gateway string `yaml:"gateway"` // network gateway
    Address string `yaml:"address"` // network address
}

type yamlConfig struct {
    Custom  []string     `yaml:"custom" json:"custom"`
    Update  asset.Config `yaml:"update" json:"update"`
    Radvd   radvd.Config `yaml:"radvd" json:"radvd"`
    Proxy   proxy.Config `yaml:"proxy" json:"proxy"`
    Network struct {
        DNS    []string      `yaml:"dns" json:"dns"`
        ByPass []string      `yaml:"bypass" json:"bypass"`
        IPv4   yamlNetConfig `yaml:"ipv4" json:"ipv4"`
        IPv6   yamlNetConfig `yaml:"ipv6" json:"ipv6"`
    } `yaml:"network" json:"network"`
}

func yamlDecode(raw []byte) yamlConfig {
    var config yamlConfig
    log.Debugf("Decode yaml content -> \n%s", string(raw))
    if err := yaml.Unmarshal(raw, &config); err != nil { // yaml (or json) decode
        log.Panicf("Decode config file error -> %v", err)
    }
    log.Debugf("Decoded config -> %v", config)
    return config
}

func decodeDns(rawConfig *yamlConfig, config *Config) {
    for _, address := range rawConfig.Network.DNS { // dns options
        if common.IsIPv4(address, false) || common.IsIPv6(address, false) {
            config.DNS = append(config.DNS, address)
        } else {
            log.Panicf("Invalid DNS server -> %s", address)
        }
    }
    log.Debugf("DNS server -> %v", config.DNS)
}

func decodeBypass(rawConfig *yamlConfig, config *Config) {
    for _, address := range rawConfig.Network.ByPass { // bypass options
        if common.IsIPv4(address, true) {
            config.IPv4.Bypass = append(config.IPv4.Bypass, address)
        } else if common.IsIPv6(address, true) {
            config.IPv6.Bypass = append(config.IPv6.Bypass, address)
        } else {
            log.Panicf("Invalid bypass CIDR -> %s", address)
        }
    }
    log.Debugf("IPv4 bypass CIDR -> %s", config.IPv4.Bypass)
    log.Debugf("IPv6 bypass CIDR -> %s", config.IPv6.Bypass)
}

func decodeIPv4(rawConfig *yamlConfig, config *Config) {
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

func decodeIPv6(rawConfig *yamlConfig, config *Config) {
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

func decodeProxy(rawConfig *yamlConfig, config *Config) {
    config.Proxy = rawConfig.Proxy
    log.Debugf("Proxy log level -> %s", config.Proxy.Log)
    log.Debugf("Http inbounds -> %v", config.Proxy.Http)
    log.Debugf("Socks5 inbounds -> %v", config.Proxy.Socks)
    log.Debugf("Add-on inbounds -> %v", config.Proxy.AddOn)
    log.Debugf("Connection sniff -> %t", config.Proxy.Sniff.Enable)
    log.Debugf("Connection redirect -> %t", config.Proxy.Sniff.Redirect)
    log.Debugf("Connection sniff exlcude -> %v", config.Proxy.Sniff.Exclude)
}

func decodeRadvd(rawConfig *yamlConfig, config *Config) {
    config.Radvd = rawConfig.Radvd
    log.Debugf("Radvd enable -> %t", config.Radvd.Enable)
    log.Debugf("Radvd options -> %v", config.Radvd.Option)
    log.Debugf("Radvd prefix -> %v", config.Radvd.Prefix)
    log.Debugf("Radvd route -> %v", config.Radvd.Route)
    log.Debugf("Radvd clients -> %v", config.Radvd.Client)
    log.Debugf("Radvd RDNSS -> %v", config.Radvd.RDNSS)
    log.Debugf("Radvd DNSSL -> %v", config.Radvd.DNSSL)
}

func decodeUpdate(rawConfig *yamlConfig, config *Config) {
    config.Update = rawConfig.Update
    log.Debugf("Update cron -> %s", config.Update.Cron)
    log.Debugf("Update urls -> %v", config.Update.Url)
}

func decodeCustom(rawConfig *yamlConfig) []string {
    customScript := rawConfig.Custom
    log.Debugf("Custom script -> %v", customScript)
    return customScript
}

func decode(rawConfig yamlConfig) Config {
    var config Config
    decodeDns(&rawConfig, &config)
    decodeBypass(&rawConfig, &config)
    decodeIPv4(&rawConfig, &config)
    decodeIPv6(&rawConfig, &config)
    decodeProxy(&rawConfig, &config)
    decodeUpdate(&rawConfig, &config)
    config.Script = decodeCustom(&rawConfig)
    decodeRadvd(&rawConfig, &config)
    return config
}
