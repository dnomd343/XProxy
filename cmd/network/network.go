package network

import (
    log "github.com/sirupsen/logrus"
    "regexp"
)

func getV4Cidr() []string {
    var v4Cidr []string
    _, output := run("ip", "-4", "addr")
    for _, temp := range regexp.MustCompile(`inet (\S+)`).FindAllStringSubmatch(output, -1) {
        v4Cidr = append(v4Cidr, temp[1])
    }
    return v4Cidr
}

func getV6Cidr() []string {
    var v6Cidr []string
    _, output := run("ip", "-6", "addr")
    for _, temp := range regexp.MustCompile(`inet6 (\S+)`).FindAllStringSubmatch(output, -1) {
        v6Cidr = append(v6Cidr, temp[1])
    }
    return v6Cidr
}

func flushNetwork(dev string) {
    log.Info("Flush system IP configure")
    run("ip", "link", "set", dev, "down")
    run("ip", "-4", "addr", "flush", "dev", dev)
    run("ip", "-6", "addr", "flush", "dev", dev)
    run("ip", "link", "set", dev, "up")
}

func flushV4Network(dev string) {
    log.Info("Flush system IPv4 configure")
    run("ip", "link", "set", dev, "down")
    run("ip", "-4", "addr", "flush", "dev", dev)
    run("ip", "link", "set", dev, "up")
}

func flushV6Network(dev string) {
    log.Info("Flush system IPv6 configure")
    run("ip", "link", "set", dev, "down")
    run("ip", "-6", "addr", "flush", "dev", dev)
    run("ip", "link", "set", dev, "up")
}

func loadV4Network(v4 *Config, dev string) {
    log.Info("Enabled IPv4 forward")
    run("sysctl", "-w", "net.ipv4.ip_forward=1")
    log.Info("Setting up system IPv4 configure")
    if v4.Address != "" {
        run("ip", "-4", "addr", "add", v4.Address, "dev", dev)
    }
    if v4.Gateway != "" {
        run("ip", "-4", "route", "add", "default", "via", v4.Gateway, "dev", dev)
    }
}

func loadV6Network(v6 *Config, dev string) {
    log.Info("Enabled IPv6 forward")
    run("sysctl", "-w", "net.ipv6.conf.all.forwarding=1")
    log.Info("Setting up system IPv6 configure")
    if v6.Address != "" {
        run("ip", "-6", "addr", "add", v6.Address, "dev", dev)
    }
    if v6.Gateway != "" {
        run("ip", "-6", "route", "add", "default", "via", v6.Gateway, "dev", dev)
    }
}
