package network

import (
    log "github.com/sirupsen/logrus"
    "regexp"
    "time"
)

func getV4Cidr() []string { // fetch ipv4 network range
    var v4Cidr []string
    _, output := run("ip", "-4", "addr")
    for _, temp := range regexp.MustCompile(`inet (\S+)`).FindAllStringSubmatch(output, -1) {
        v4Cidr = append(v4Cidr, temp[1])
    }
    return v4Cidr
}

func getV6Cidr() []string { // fetch ipv6 network range
    var v6Cidr []string
    _, output := run("ip", "-6", "addr")
    for _, temp := range regexp.MustCompile(`inet6 (\S+)`).FindAllStringSubmatch(output, -1) {
        v6Cidr = append(v6Cidr, temp[1])
    }
    return v6Cidr
}

func enableIpForward() { // enable ip forward function
    log.Info("Enabled IPv4 forward")
    run("sysctl", "-w", "net.ipv4.ip_forward=1")
    log.Info("Enabled IPv6 forward")
    run("sysctl", "-w", "net.ipv6.conf.all.forwarding=1")
}

func flushNetwork(dev string, isV4 bool, isV6 bool) { // flush ipv4 and ipv6 network
    log.Info("Flush system IP configure")
    run("ip", "link", "set", dev, "down")
    if isV4 {
        run("ip", "-4", "addr", "flush", "dev", dev)
    }
    if isV6 {
        run("ip", "-6", "addr", "flush", "dev", dev)
    }
    run("ip", "link", "set", dev, "up")
}

func loadV4Network(v4 *Config, dev string) { // setting up ipv4 network
    log.Info("Setting up system IPv4 configure")
    if v4.Address != "" {
        run("ip", "-4", "addr", "add", v4.Address, "dev", dev)
    }
    if v4.Gateway != "" {
        run("ip", "-4", "route", "add", "default", "via", v4.Gateway, "dev", dev)
    }
}

func loadV6Network(v6 *Config, dev string) { // setting up ipv6 network
    log.Info("Setting up system IPv6 configure")
    if v6.Address != "" {
        run("ip", "-6", "addr", "add", v6.Address, "dev", dev)
    }
    if v6.Gateway != "" {
        run("ip", "-6", "route", "add", "default", "via", v6.Gateway, "dev", dev)
    }
}

func loadNetwork(dev string, v4 *Config, v6 *Config) {
    setV4 := v4.Address != "" || v4.Gateway != ""
    setV6 := v6.Address != "" || v6.Gateway != ""
    if setV4 && setV6 { // load both ipv4 and ipv6
        flushNetwork(dev, true, true)
        loadV4Network(v4, dev)
        loadV6Network(v6, dev)
    } else if setV4 { // only load ipv4 network
        flushNetwork(dev, true, false)
        loadV4Network(v4, dev)
    } else if setV6 { // only load ipv6 network
        flushNetwork(dev, false, true)
        loadV6Network(v6, dev)
    } else { // skip network settings
        log.Infof("Skip system IP configure")
    }
    if setV6 {
        log.Info("Wait 1s for IPv6 setting up")
        time.Sleep(time.Second) // wait for ipv6 setting up (ND protocol) -> RA should reply less than 0.5s
    }
}
