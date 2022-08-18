package network

import (
    "XProxy/cmd/common"
    log "github.com/sirupsen/logrus"
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
    log.Info("Wait 1s for IPv6 setting up")
    time.Sleep(time.Second) // wait for ipv6 setting up (ND protocol) -> RA should reply less than 0.5s
    loadV4TProxy(ipv4, getV4Cidr())
    loadV6TProxy(ipv6, getV6Cidr())
}
