package main

import (
    log "github.com/sirupsen/logrus"
    "os"
    "strconv"
)

func loadDns() {
    if len(dnsServer) == 0 {
        log.Info("Using system DNS server")
        return
    }
    log.Infof("Setting up DNS server -> %v", dnsServer)
    dnsContent := ""
    for _, address := range dnsServer {
        dnsContent += "nameserver " + address + "\n"
    }
    err := os.WriteFile("/etc/resolv.conf", []byte(dnsContent), 0644)
    if err != nil {
        log.Error("Setting up DNS failed")
        panic("Setting up DNS failed")
    }
}

func loadNetwork() {
    log.Info("Enabled IP forward")
    runCommand([]string{"sysctl", "-w", "net.ipv4.ip_forward=1"})
    runCommand([]string{"sysctl", "-w", "net.ipv6.conf.all.forwarding=1"})

    log.Info("Flush system IP configure")
    runCommand([]string{"ip", "link", "set", "eth0", "down"})
    runCommand([]string{"ip", "-4", "addr", "flush", "dev", "eth0"})
    runCommand([]string{"ip", "-6", "addr", "flush", "dev", "eth0"})
    runCommand([]string{"ip", "link", "set", "eth0", "down"})

    log.Info("Setting up system IP configure")
    if v4Address != "" {
        runCommand([]string{"ip", "-4", "addr", "add", v4Address, "dev", "eth0"})
    }
    if v4Gateway != "" {
        runCommand([]string{"ip", "-4", "route", "add", "default", "via", v4Gateway})
    }
    if v6Address != "" {
        runCommand([]string{"ip", "-6", "addr", "add", v6Address, "dev", "eth0"})
    }
    if v6Gateway != "" {
        runCommand([]string{"ip", "-6", "route", "add", "default", "via", v6Gateway})
    }
}

func loadTProxy() {
    log.Info("Setting up TProxy of IPv4")
    v4TableNum := strconv.Itoa(v4RouteTable)
    runCommand([]string{"ip", "-4", "rule", "add", "fwmark", "1", "table", v4TableNum})
    runCommand([]string{"ip", "-4", "route", "add", "local", "0.0.0.0/0", "dev", "lo", "table", v4TableNum})
    runCommand([]string{"iptables", "-t", "mangle", "-N", "XPROXY"})
    for _, cidr := range v4Bypass {
        runCommand([]string{"iptables", "-t", "mangle", "-A", "XPROXY", "-d", cidr, "-j", "RETURN"})
    }
    runCommand([]string{"iptables", "-t", "mangle", "-A", "XPROXY", "-p", "tcp", "-j", "TPROXY",
        "--on-port", strconv.Itoa(v4TProxyPort), "--tproxy-mark", "1"})
    runCommand([]string{"iptables", "-t", "mangle", "-A", "XPROXY", "-p", "udp", "-j", "TPROXY",
        "--on-port", strconv.Itoa(v4TProxyPort), "--tproxy-mark", "1"})
    runCommand([]string{"iptables", "-t", "mangle", "-A", "PREROUTING", "-j", "XPROXY"})

    log.Info("Setting up TProxy of IPv6")
    v6TableNum := strconv.Itoa(v6RouteTable)
    runCommand([]string{"ip", "-6", "rule", "add", "fwmark", "1", "table", v6TableNum})
    runCommand([]string{"ip", "-6", "route", "add", "local", "::/0", "dev", "lo", "table", v6TableNum})
    runCommand([]string{"ip6tables", "-t", "mangle", "-N", "XPROXY6"})
    for _, cidr := range v6Bypass {
        runCommand([]string{"ip6tables", "-t", "mangle", "-A", "XPROXY6", "-d", cidr, "-j", "RETURN"})
    }
    runCommand([]string{"ip6tables", "-t", "mangle", "-A", "XPROXY6", "-p", "tcp", "-j", "TPROXY",
        "--on-port", strconv.Itoa(v6TProxyPort), "--tproxy-mark", "1"})
    runCommand([]string{"ip6tables", "-t", "mangle", "-A", "XPROXY6", "-p", "udp", "-j", "TPROXY",
        "--on-port", strconv.Itoa(v6TProxyPort), "--tproxy-mark", "1"})
    runCommand([]string{"ip6tables", "-t", "mangle", "-A", "PREROUTING", "-j", "XPROXY6"})
}
