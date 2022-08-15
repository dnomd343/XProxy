package main

import (
    "fmt"
    "gopkg.in/yaml.v3"
    "net"
    "strconv"
    "strings"
)

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
    fmt.Println(config)
}
