package main

import (
    "gopkg.in/yaml.v3"
    "net"
    "strconv"
    "strings"
)

var v4Bypass []string
var v6Bypass []string
var dnsServer []string

var v4Gateway string
var v4Address string
var v4Forward bool

var v6Gateway string
var v6Address string
var v6Forward bool

type netConfig struct {
    Gateway string `yaml:"gateway"` // network gateway
    Address string `yaml:"address"` // network address
    Forward bool   `yaml:"forward"` // enabled net forward
}

type Config struct {
    Network struct {
        DNS    []string  `yaml:"dns"`    // system dns server
        ByPass []string  `yaml:"bypass"` // cidr bypass list
        IPv4   netConfig `yaml:"ipv4"`   // ipv4 network configure
        IPv6   netConfig `yaml:"ipv6"`   // ipv6 network configure
    }
}

func isIP(ipAddr string, isRange bool, ipLength int, ipFlag string) bool {
    var address string
    if isRange {
        temp := strings.Split(ipAddr, "/")
        if len(temp) != 2 { // not {IP_ADDRESS}/{LENGTH} format
            return false
        }
        length, err := strconv.Atoi(temp[1])
        if err != nil { // range length not a integer
            return false
        }
        if length < 0 || length > ipLength { // length should between 0 ~ ipLength
            return false
        }
        address = temp[0]
    } else {
        address = ipAddr
    }
    ip := net.ParseIP(address) // try to convert ip
    return ip != nil && strings.Contains(address, ipFlag)
}

func isIPv4(ipAddr string, isRange bool) bool {
    return isIP(ipAddr, isRange, 32, ".")
}

func isIPv6(ipAddr string, isRange bool) bool {
    return isIP(ipAddr, isRange, 128, ":")
}

func loadConfig(rawConfig []byte) {
    config := Config{}
    err := yaml.Unmarshal(rawConfig, &config) // yaml (or json) decode
    if err != nil {
        panic(err)
    }

    for _, address := range config.Network.DNS { // dns options
        if isIPv4(address, false) || isIPv6(address, false) {
            dnsServer = append(dnsServer, address)
        } else {
            panic("Invalid DNS server -> " + address)
        }
    }

    for _, address := range config.Network.ByPass { // bypass options
        if isIPv4(address, true) {
            v4Bypass = append(v4Bypass, address)
        } else if isIPv6(address, true) {
            v6Bypass = append(v6Bypass, address)
        } else {
            panic("Invalid bypass CIDR -> " + address)
        }
    }

    v4Address = config.Network.IPv4.Address // ipv4 options
    v4Gateway = config.Network.IPv4.Gateway
    v4Forward = config.Network.IPv4.Forward
    if v4Address != "" && !isIPv4(v4Address, true) {
        panic("Invalid IPv4 address -> " + v4Address)
    }
    if v4Gateway != "" && !isIPv4(v4Gateway, false) {
        panic("Invalid IPv4 gateway -> " + v4Gateway)
    }

    v6Address = config.Network.IPv6.Address // ipv6 options
    v6Gateway = config.Network.IPv6.Gateway
    v6Forward = config.Network.IPv6.Forward
    if v6Address != "" && !isIPv6(v6Address, true) {
        panic("Invalid IPv6 address -> " + v6Address)
    }
    if v6Gateway != "" && !isIPv6(v6Gateway, false) {
        panic("Invalid IPv6 gateway -> " + v6Gateway)
    }
}
