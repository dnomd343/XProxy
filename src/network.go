package main

import (
    log "github.com/sirupsen/logrus"
    "os"
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
