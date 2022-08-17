package main

import (
    log "github.com/sirupsen/logrus"
)

var v4Bypass []string
var v6Bypass []string
var dnsServer []string

var v4Gateway string
var v4Address string
var v6Gateway string
var v6Address string

func loadConfig(configFile string) {

    config := Config{}

    enableSniff = config.Proxy.Sniff
    log.Infof("Connection sniff -> %v", enableSniff)
    enableRedirect = config.Proxy.Redirect
    log.Infof("Connection redirect -> %v", enableRedirect)
    httpInbounds = config.Proxy.Http
    log.Infof("Http inbounds -> %v", httpInbounds)
    socksInbounds = config.Proxy.Socks
    log.Infof("Socks5 inbounds -> %v", socksInbounds)
    addOnInbounds = config.Proxy.AddOn
    log.Infof("Add-on inbounds -> %v", addOnInbounds)

    updateCron = config.Update.Cron
    log.Infof("Update cron -> %s", updateCron)
    updateUrls = config.Update.Url
    log.Infof("Update url -> %v", updateUrls)

    preScript = config.Script
    log.Infof("Pre-script -> %v", preScript)
}
