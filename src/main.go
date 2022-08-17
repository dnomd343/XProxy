package main

import (
    log "github.com/sirupsen/logrus"
    "os"
    "os/signal"
    "syscall"
)

var preScript []string

func main() {
    defer func() {
        if err := recover(); err != nil {
            log.Errorf("Unknown error -> %v", err)
        }
    }()

    for _, script := range preScript {
        log.Infof("Run script command -> %s", script)
        runCommand("sh", "-c", script)
    }

    xray := newProcess("xray", "-confdir", "/etc/xproxy/config")
    xray.startProcess(true, true)
    subProcess = append(subProcess, xray)
    daemon()

    sigExit := make(chan os.Signal, 1)
    signal.Notify(sigExit, syscall.SIGINT, syscall.SIGTERM)
    <-sigExit
    exit()
}
