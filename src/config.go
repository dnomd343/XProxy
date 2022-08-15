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

func isIP(ipAddr string, isRange bool, allowEmpty bool, ipLength int, ipFlag string) bool {
    var address string
    if allowEmpty && ipAddr == "" { // empty case
        return true
    }
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

func isIPv4(ipAddr string, isRange bool, allowEmpty bool) bool {
    return isIP(ipAddr, isRange, allowEmpty, 32, ".")
}

func isIPv6(ipAddr string, isRange bool, allowEmpty bool) bool {
    return isIP(ipAddr, isRange, allowEmpty, 128, ":")
}

func loadConfig(rawConfig []byte) {
    config := Config{}
    err := yaml.Unmarshal(rawConfig, &config) // yaml (or json) decode
    if err != nil {
        panic(err)
    }
    for _, address := range config.Network.DNS { // load dns configure
        if isIPv4(address, false, false) || isIPv6(address, false, false) {
            dnsServer = append(dnsServer, address)
        } else {
            panic("Invalid DNS server -> " + address)
        }
    }
    for _, address := range config.Network.ByPass { // load bypass configure
        if isIPv4(address, true, false) {
            v4Bypass = append(v4Bypass, address)
        } else if isIPv6(address, true, false) {
            v6Bypass = append(v6Bypass, address)
        } else {
            panic("Invalid bypass CIDR -> " + address)
        }
    }
    //fmt.Println(config)
}
