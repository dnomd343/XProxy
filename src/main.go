package main

import (
    log "github.com/sirupsen/logrus"
    "os"
)

func main() {
    log.SetLevel(log.DebugLevel)
    log.Warning("XProxy start")
    content, err := os.ReadFile("test.yml")
    if err != nil {
        panic(err)
    }
    loadConfig(content)
}
