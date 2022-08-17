package network

import (
    "XProxy/cmd/common"
    "time"
)

type Config struct {
    RouteTable int
    TProxyPort int
    Address    string
    Gateway    string
    Bypass     []string
}

var run = common.RunCommand

func Load(dns []string, ipv4 Config, ipv6 Config) {
    loadDns(dns)   // init dns server
    flushNetwork() // clear network settings
    loadV4Network(ipv4)
    loadV6Network(ipv6)
    time.Sleep(time.Second) // wait 1s for ipv6 (ND protocol)
    loadV4TProxy(ipv4, getV4Cidr())
    loadV6TProxy(ipv6, getV6Cidr())
}
