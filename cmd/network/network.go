package network

import (
    "XProxy/cmd/common"
    log "github.com/sirupsen/logrus"
    "regexp"
)

func getV4Cidr() []string {
    var v4Cidr []string
    _, output := common.RunCommand("ip", "-4", "addr")
    for _, temp := range regexp.MustCompile(`inet (\S+)`).FindAllStringSubmatch(output, -1) {
        v4Cidr = append(v4Cidr, temp[1])
    }
    return v4Cidr
}

func getV6Cidr() []string {
    var v6Cidr []string
    _, output := common.RunCommand("ip", "-6", "addr")
    for _, temp := range regexp.MustCompile(`inet6 (\S+)`).FindAllStringSubmatch(output, -1) {
        v6Cidr = append(v6Cidr, temp[1])
    }
    return v6Cidr
}

func flushNetwork() {
    log.Info("Flush system IP configure")
    common.RunCommand("ip", "link", "set", "eth0", "down")
    common.RunCommand("ip", "-4", "addr", "flush", "dev", "eth0")
    common.RunCommand("ip", "-6", "addr", "flush", "dev", "eth0")
    common.RunCommand("ip", "link", "set", "eth0", "down")
}

func loadV4Network(v4 Config) {
    log.Info("Enabled IPv4 forward")
    common.RunCommand("sysctl", "-w", "net.ipv4.ip_forward=1")
    log.Info("Setting up system IPv4 configure")
    if v4.Address != "" {
        common.RunCommand("ip", "-4", "addr", "add", v4.Address, "dev", "eth0")
    }
    if v4.Gateway != "" {
        common.RunCommand("ip", "-4", "route", "add", "default", "via", v4.Gateway)
    }
}

func loadV6Network(v6 Config) {
    log.Info("Enabled IPv6 forward")
    common.RunCommand("sysctl", "-w", "net.ipv6.conf.all.forwarding=1")
    log.Info("Setting up system IPv6 configure")
    if v6.Address != "" {
        common.RunCommand("ip", "-6", "addr", "add", v6.Address, "dev", "eth0")
    }
    if v6.Gateway != "" {
        common.RunCommand("ip", "-6", "route", "add", "default", "via", v6.Gateway)
    }
}
