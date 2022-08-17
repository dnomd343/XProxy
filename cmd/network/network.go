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

func loadNetwork(v4Address string, v4Gateway string, v6Address string, v6Gateway string) {
    log.Info("Enabled IP forward")
    common.RunCommand("sysctl", "-w", "net.ipv4.ip_forward=1")
    common.RunCommand("sysctl", "-w", "net.ipv6.conf.all.forwarding=1")

    log.Info("Flush system IP configure")
    common.RunCommand("ip", "link", "set", "eth0", "down")
    common.RunCommand("ip", "-4", "addr", "flush", "dev", "eth0")
    common.RunCommand("ip", "-6", "addr", "flush", "dev", "eth0")
    common.RunCommand("ip", "link", "set", "eth0", "down")

    log.Info("Setting up system IP configure")
    if v4Address != "" {
        common.RunCommand("ip", "-4", "addr", "add", v4Address, "dev", "eth0")
    }
    if v4Gateway != "" {
        common.RunCommand("ip", "-4", "route", "add", "default", "via", v4Gateway)
    }
    if v6Address != "" {
        common.RunCommand("ip", "-6", "addr", "add", v6Address, "dev", "eth0")
    }
    if v6Gateway != "" {
        common.RunCommand("ip", "-6", "route", "add", "default", "via", v6Gateway)
    }
}
