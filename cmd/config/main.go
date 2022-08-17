package config

import (
    "XProxy/cmd/common"
    log "github.com/sirupsen/logrus"
    "os"
)

type Config struct {
    // asset update

    DNS []string

    V4Address string
    V4Gateway string
    V4Bypass  []string

    V6Address string
    V6Gateway string
    V6Bypass  []string

    // httpInbounds
    // socksInbounds
    // addOnInbounds

}

func Load(configFile string) Config {
    if !common.IsFileExist(configFile) { // configure not exist -> load default
        log.Infof("Load default configure -> %s", configFile)
        common.WriteFile(configFile, defaultConfig, false)
    }
    raw, err := os.ReadFile(configFile) // read configure content
    if err != nil {
        log.Panicf("Failed to open %s -> %v", configFile, err)
    }
    rawConfig := yamlDecode(raw) // decode yaml content
    return decode(rawConfig)
}
