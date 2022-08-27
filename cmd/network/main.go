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
    Exclude    []string
}

var run = common.RunCommand

func Load(dns []string, ipv4 *Config, ipv6 *Config, dev string) {
    loadDns(dns) // init dns server
    delay := false
    setV4 := ipv4.Address != "" || ipv4.Gateway != ""
    setV6 := ipv6.Address != "" || ipv6.Gateway != ""
    if setV4 && setV6 { // clear network settings
        delay = true
        flushNetwork(dev)
        loadV4Network(ipv4, dev)
        loadV6Network(ipv6, dev)
    } else if setV6 {
        delay = true
        flushV6Network(dev)
        loadV6Network(ipv6, dev)
    } else if setV4 {
        flushV4Network(dev)
        loadV4Network(ipv4, dev)
    } else {
        log.Infof("Skip system IP configure")
    }
    if delay {
        log.Info("Wait 1s for IPv6 setting up")
        time.Sleep(time.Second) // wait for ipv6 setting up (ND protocol) -> RA should reply less than 0.5s
    }
    loadV4TProxy(ipv4, getV4Cidr())
    loadV6TProxy(ipv6, getV6Cidr())
}
