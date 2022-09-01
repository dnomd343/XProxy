package process

import (
    log "github.com/sirupsen/logrus"
    "time"
)

func daemonSub(sub *Process) {
    for sub.process.ProcessState == nil { // until process exit
        sub.Wait()
    }
    log.Warningf("Catch process %s exit", sub.name)
    time.Sleep(5 * time.Second) // delay 3s -> try to restart
    if !exitFlag {
        sub.Run(true, sub.env)
        log.Infof("Process %s restart success", sub.name)
        daemonSub(sub)
    }
}

func (p *Process) Daemon() {
    if p.process == nil { // process not running
        log.Infof("Process %s disabled -> skip daemon", p.name)
        return
    }
    log.Infof("Daemon of process %s start", p.name)
    go func() {
        daemonSub(p) // start daemon process
        log.Infof("Process %s exit daemon mode", p.name)
    }()
}
