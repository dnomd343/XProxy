package network

import (
	"XProxy/cmd/common"
	log "github.com/sirupsen/logrus"
)

func loadDns(dns []string) {
	if len(dns) == 0 { // without dns server
		log.Info("Using system DNS server")
		return
	}
	log.Infof("Setting up DNS server -> %v", dns)
	dnsConfig := ""
	for _, dnsAddr := range dns {
		dnsConfig += "nameserver " + dnsAddr + "\n"
	}
	common.WriteFile("/etc/resolv.conf", dnsConfig, true)
}
