package process

import (
    log "github.com/sirupsen/logrus"
    "syscall"
)

var exitFlag bool

func Exit(subProcess ...*Process) {
    exitFlag = true // setting up exit flag -> exit daemon mode
    log.Warningf("Start exit process")
    for _, sub := range subProcess {
        if sub.process != nil {
            log.Infof("Send kill signal to process %s", sub.name)
            sub.Signal(syscall.SIGTERM)
        }
    }
    log.Info("Wait all sub process exit")
    for _, sub := range subProcess {
        if sub.process != nil {
            _ = sub.process.Wait()
        }
    }
    log.Infof("Exit complete")
}
