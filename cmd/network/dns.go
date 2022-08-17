package network

import (
    "XProxy/cmd/common"
    log "github.com/sirupsen/logrus"
)

func loadDns(dnsServer []string) {
    if len(dnsServer) == 0 {
        log.Info("Using system DNS server")
        return
    }
    log.Infof("Setting up DNS server -> %v", dnsServer)
    dnsConfig := ""
    for _, address := range dnsServer {
        dnsConfig += "nameserver " + address + "\n"
    }
    common.WriteFile("/etc/resolv.conf", dnsConfig, true)
}
