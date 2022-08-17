package proxy

import (
    "XProxy/cmd/common"
    "encoding/json"
    log "github.com/sirupsen/logrus"
)

func saveConfig(configDir string, caption string, content string, overwrite bool) {
    filePath := configDir + "/" + caption + ".json"
    common.WriteFile(filePath, content+"\n", overwrite)
}

func jsonEncode(raw interface{}) string {
    jsonOutput, _ := json.MarshalIndent(raw, "", "  ") // json encode
    return string(jsonOutput)
}

func loadLog(logLevel string, logDir string) string {
    if logLevel != "debug" && logLevel != "info" &&
        logLevel != "warning" && logLevel != "error" && logLevel != "none" {
        log.Warningf("Unknown log level -> %s", logLevel)
        logLevel = "warning" // using `warning` as default
    }
    return jsonEncode(logObject{
        Loglevel: logLevel,
        Access:   logDir + "/access.log",
        Error:    logDir + "/error.log",
    })
}

func loadInbounds(config Config) string {
    inbounds := inboundsObject{}
    sniff := sniffObject{
        Enabled:      config.Sniff,
        RouteOnly:    !config.Redirect,
        DestOverride: []string{"http", "tls"},
    }
    inbounds.Inbounds = append(inbounds.Inbounds, tproxyConfig("tproxy", config.V4TProxyPort, sniff))
    inbounds.Inbounds = append(inbounds.Inbounds, tproxyConfig("tproxy6", config.V6TProxyPort, sniff))
    for tag, port := range config.HttpInbounds {
        inbounds.Inbounds = append(inbounds.Inbounds, httpConfig(tag, port, sniff))
    }
    for tag, port := range config.SocksInbounds {
        inbounds.Inbounds = append(inbounds.Inbounds, socksConfig(tag, port, sniff))
    }
    for _, addon := range config.AddOnInbounds {
        inbounds.Inbounds = append(inbounds.Inbounds, addon)
    }
    return jsonEncode(inbounds)
}

func Load(configDir string, exposeDir string, config Config) {
    common.CreateFolder(exposeDir + "/log")
    common.CreateFolder(exposeDir + "/config")
    common.CreateFolder(configDir)
    saveConfig(exposeDir+"/config", "dns", dnsConfig, false)
    saveConfig(exposeDir+"/config", "route", routeConfig, false)
    saveConfig(exposeDir+"/config", "outbounds", outboundsConfig, false)
    saveConfig(configDir, "inbounds", loadInbounds(config), true)
    saveConfig(configDir, "log", loadLog(config.LogLevel, exposeDir+"/log"), true)
    for _, configFile := range common.ListFiles(exposeDir+"/config", ".json") {
        common.CopyFile(exposeDir+"/config/"+configFile, configDir+"/"+configFile)
    }
}
