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
    time.Sleep(100 * time.Millisecond) // delay 100ms
    if !exitFlag {                     // not in exit process
        sub.Run(true)
        log.Infof("Process %s restart success", sub.name)
        daemonSub(sub)
    }
}

func (p *Process) Daemon() {
    if p.process == nil { // process not running
        log.Infof("Process %s disabled -> skip daemon", p.name)
        return
    }
    log.Infof("Start daemon of process %s", p.name)
    go func() {
        daemonSub(p) // start daemon process
        log.Infof("Process %s exit daemon mode", p.name)
    }()
}
