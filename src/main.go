package main

import (
    log "github.com/sirupsen/logrus"
)

func main() {
    log.SetLevel(log.DebugLevel)
    log.Warning("XProxy start")

    proxyConfig("/etc/xproxy/config", "debug", "/xproxy/log")

    //content, err := os.ReadFile("test.yml")
    //if err != nil {
    //	panic(err)
    //}
    //loadConfig(content)
}
