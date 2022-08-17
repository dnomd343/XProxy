package main

import (
    "fmt"
    log "github.com/sirupsen/logrus"
)

func main() {
    log.SetLevel(log.DebugLevel)

    fmt.Println("xproxy start")

}
