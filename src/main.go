package main

import (
    log "github.com/sirupsen/logrus"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    defer func() {
        if err := recover(); err != nil {
            log.Errorf("Unknown error -> %v", err)
        }
    }()

    xray := newProcess("xray", "-confdir", "/etc/xproxy/config")
    xray.startProcess(true, true)
    subProcess = append(subProcess, xray)
    daemon()

    sigExit := make(chan os.Signal, 1)
    signal.Notify(sigExit, syscall.SIGINT, syscall.SIGTERM)
    <-sigExit
    exit()
}
