package network

import (
    "XProxy/cmd/common"
    log "github.com/sirupsen/logrus"
    "strconv"
)

type tproxyConfig struct {
    routeTable int
    tproxyPort int
    bypassCidr []string
}

func loadV4TProxy(config tproxyConfig) {
    log.Info("Setting up TProxy of IPv4")
    tableNum := strconv.Itoa(config.routeTable)
    common.RunCommand("ip", "-4", "rule", "add", "fwmark", "1", "table", tableNum)
    common.RunCommand("ip", "-4", "route", "add", "local", "0.0.0.0/0", "dev", "lo", "table", tableNum)
    common.RunCommand("iptables", "-t", "mangle", "-N", "XPROXY")
    log.Infof("Setting up IPv4 bypass CIDR -> %v", config.bypassCidr)
    for _, cidr := range config.bypassCidr {
        common.RunCommand("iptables", "-t", "mangle", "-A", "XPROXY", "-d", cidr, "-j", "RETURN")
    }
    common.RunCommand("iptables", "-t", "mangle", "-A", "XPROXY",
        "-p", "tcp", "-j", "TPROXY", "--on-port", strconv.Itoa(config.tproxyPort), "--tproxy-mark", "1")
    common.RunCommand("iptables", "-t", "mangle", "-A", "XPROXY",
        "-p", "udp", "-j", "TPROXY", "--on-port", strconv.Itoa(config.tproxyPort), "--tproxy-mark", "1")
    common.RunCommand("iptables", "-t", "mangle", "-A", "PREROUTING", "-j", "XPROXY")
}

func loadV6TProxy(config tproxyConfig) {
    log.Info("Setting up TProxy of IPv6")
    tableNum := strconv.Itoa(config.routeTable)
    common.RunCommand("ip", "-6", "rule", "add", "fwmark", "1", "table", tableNum)
    common.RunCommand("ip", "-6", "route", "add", "local", "::/0", "dev", "lo", "table", tableNum)
    common.RunCommand("ip6tables", "-t", "mangle", "-N", "XPROXY6")
    log.Infof("Setting up IPv6 bypass CIDR -> %v", config.bypassCidr)
    for _, cidr := range config.bypassCidr {
        common.RunCommand("ip6tables", "-t", "mangle", "-A", "XPROXY6", "-d", cidr, "-j", "RETURN")
    }
    common.RunCommand("ip6tables", "-t", "mangle", "-A", "XPROXY6",
        "-p", "tcp", "-j", "TPROXY", "--on-port", strconv.Itoa(config.tproxyPort), "--tproxy-mark", "1")
    common.RunCommand("ip6tables", "-t", "mangle", "-A", "XPROXY6",
        "-p", "udp", "-j", "TPROXY", "--on-port", strconv.Itoa(config.tproxyPort), "--tproxy-mark", "1")
    common.RunCommand("ip6tables", "-t", "mangle", "-A", "PREROUTING", "-j", "XPROXY6")
}
