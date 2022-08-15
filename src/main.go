package main

import (
    log "github.com/sirupsen/logrus"
)

var logLevel = "debug"

func main() {
    log.SetLevel(log.DebugLevel)
    log.Warning("XProxy start")

    loadProxy("/etc/xproxy/config", "/xproxy")

    //content, err := os.ReadFile("test.yml")
    //if err != nil {
    //	panic(err)
    //}
    //loadConfig(content)
}
