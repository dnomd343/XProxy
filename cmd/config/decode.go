package config

import (
    "XProxy/cmd/asset"
    "XProxy/cmd/common"
    "XProxy/cmd/radvd"
    log "github.com/sirupsen/logrus"
    "gopkg.in/yaml.v3"
)

type yamlNetConfig struct {
    Gateway string `yaml:"gateway"` // network gateway
    Address string `yaml:"address"` // network address
}

type yamlConfig struct {
    Custom []string     `yaml:"custom"`
    Update asset.Config `yaml:"update"`
    Proxy  struct {
        Log   string `yaml:"log"`
        Sniff struct {
            Enable   bool     `yaml:"enable"`
            Redirect bool     `yaml:"redirect"`
            Exclude  []string `yaml:"exclude"`
        } `yaml:"sniff"`
        Http  map[string]int `yaml:"http"`
        Socks map[string]int `yaml:"socks"`
        AddOn []interface{}  `yaml:"addon"`
    } `yaml:"proxy"`
    Network struct {
        DNS    []string      `yaml:"dns"`    // system dns server
        ByPass []string      `yaml:"bypass"` // cidr bypass list
        IPv4   yamlNetConfig `yaml:"ipv4"`   // ipv4 network configure
        IPv6   yamlNetConfig `yaml:"ipv6"`   // ipv6 network configure
    } `yaml:"network"`
    Radvd radvd.Config `yaml:"radvd"`
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

func decodeDns(rawConfig *yamlConfig) []string {
    var dns []string
    for _, address := range rawConfig.Network.DNS { // dns options
        if common.IsIPv4(address, false) || common.IsIPv6(address, false) {
            dns = append(dns, address)
        } else {
            log.Panicf("Invalid DNS server -> %s", address)
        }
    }
    log.Debugf("DNS server -> %v", dns)
    return dns
}

func decodeBypass(rawConfig *yamlConfig) ([]string, []string) {
    var v4Bypass, v6Bypass []string
    for _, address := range rawConfig.Network.ByPass { // bypass options
        if common.IsIPv4(address, true) {
            v4Bypass = append(v4Bypass, address)
        } else if common.IsIPv6(address, true) {
            v6Bypass = append(v6Bypass, address)
        } else {
            log.Panicf("Invalid bypass CIDR -> %s", address)
        }
    }
    log.Debugf("IPv4 bypass CIDR -> %s", v4Bypass)
    log.Debugf("IPv6 bypass CIDR -> %s", v6Bypass)
    return v4Bypass, v6Bypass
}

func decodeIPv4(rawConfig *yamlConfig) (string, string) {
    v4Address := rawConfig.Network.IPv4.Address
    v4Gateway := rawConfig.Network.IPv4.Gateway
    if v4Address != "" && !common.IsIPv4(v4Address, true) {
        log.Panicf("Invalid IPv4 address -> %s", v4Address)
    }
    if v4Gateway != "" && !common.IsIPv4(v4Gateway, false) {
        log.Panicf("Invalid IPv4 gateway -> %s", v4Gateway)
    }
    log.Debugf("IPv4 -> address = %s | gateway = %s", v4Address, v4Gateway)
    return v4Address, v4Gateway
}

func decodeIPv6(rawConfig *yamlConfig) (string, string) {
    v6Address := rawConfig.Network.IPv6.Address
    v6Gateway := rawConfig.Network.IPv6.Gateway
    if v6Address != "" && !common.IsIPv6(v6Address, true) {
        log.Panicf("Invalid IPv6 address -> %s", v6Address)
    }
    if v6Gateway != "" && !common.IsIPv6(v6Gateway, false) {
        log.Panicf("Invalid IPv6 gateway -> %s", v6Gateway)
    }
    log.Debugf("IPv6 -> address = %s | gateway = %s", v6Address, v6Gateway)
    return v6Address, v6Gateway
}

func decodeProxy(rawConfig *yamlConfig, config *Config) {
    config.EnableSniff = rawConfig.Proxy.Sniff.Enable
    log.Debugf("Connection sniff -> %t", config.EnableSniff)
    config.EnableRedirect = rawConfig.Proxy.Sniff.Redirect
    log.Debugf("Connection redirect -> %t", config.EnableRedirect)
    config.SniffExclude = rawConfig.Proxy.Sniff.Exclude
    log.Debugf("Connection sniff exlcude -> %v", config.SniffExclude)
    config.HttpInbounds = rawConfig.Proxy.Http
    log.Debugf("Http inbounds -> %v", config.HttpInbounds)
    config.SocksInbounds = rawConfig.Proxy.Socks
    log.Debugf("Socks5 inbounds -> %v", config.SocksInbounds)
    config.AddOnInbounds = rawConfig.Proxy.AddOn
    log.Debugf("Add-on inbounds -> %v", config.AddOnInbounds)
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
    config.LogLevel = rawConfig.Proxy.Log
    config.DNS = decodeDns(&rawConfig)
    config.V4Bypass, config.V6Bypass = decodeBypass(&rawConfig)
    config.V4Address, config.V4Gateway = decodeIPv4(&rawConfig)
    config.V6Address, config.V6Gateway = decodeIPv6(&rawConfig)
    decodeProxy(&rawConfig, &config)
    decodeUpdate(&rawConfig, &config)
    config.Script = decodeCustom(&rawConfig)
    decodeRadvd(&rawConfig, &config)
    return config
}
