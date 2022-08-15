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
