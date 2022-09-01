package network

import (
    "XProxy/cmd/common"
)

type Config struct {
    RouteTable int
    TProxyPort int
    Address    string
    Gateway    string
    Bypass     []string
    Exclude    []string
}

var run = common.RunCommand

func Load(dns []string, dev string, ipv4 *Config, ipv6 *Config) {
    loadDns(dns) // init dns server
    enableIpForward()
    loadNetwork(dev, ipv4, ipv6)
    loadV4TProxy(ipv4, getV4Cidr())
    loadV6TProxy(ipv6, getV6Cidr())
}
