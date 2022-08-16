package main

import (
    log "github.com/sirupsen/logrus"
    "os"
    "os/exec"
    "syscall"
    "time"
)

var exitFlag bool
var subProcess []*Process

type Process struct {
    name    string
    command []string
    process *exec.Cmd
}

func newProcess(command ...string) *Process {
    process := new(Process)
    process.name = command[0]
    process.command = command
    log.Debugf("New process %s -> %v", process.name, process.command)
    return process
}

func (p *Process) startProcess(isStdout bool, isStderr bool) {
    p.process = exec.Command(p.command[0], p.command[1:]...)
    if isStdout {
        p.process.Stdout = os.Stdout
    }
    if isStderr {
        p.process.Stderr = os.Stderr
    }
    err := p.process.Start()
    if err != nil {
        log.Errorf("Failed to start %s -> %v", p.name, err)
    }
    log.Infof("Start process %s -> PID = %d", p.name, p.process.Process.Pid)
}

func (p *Process) sendSignal(signal syscall.Signal) {
    if p.process != nil {
        log.Debugf("Send signal %v to %s", signal, p.name)
        _ = p.process.Process.Signal(signal)
    }
}

func (p *Process) waitProcess() {
    if p.process != nil {
        err := p.process.Wait()
        if err != nil {
            log.Warningf("Wait process %s -> %v", p.name, err)
        }
    }
}

func daemonSub(sub *Process) {
    for sub.process.ProcessState == nil {
        sub.waitProcess()
    }
    log.Warningf("Catch process %s exit", sub.name)
    time.Sleep(10 * time.Millisecond) // delay 10ms
    if !exitFlag {
        sub.startProcess(true, true)
        log.Infof("Process %s restart success", sub.name)
        daemonSub(sub)
    }
}

func daemon() {
    for _, sub := range subProcess {
        if sub.process == nil {
            log.Infof("Process %s disabled -> skip daemon", sub.name)
            return
        }
        log.Infof("Start daemon of process %s", sub.name)
        sub := sub
        go func() {
            daemonSub(sub)
            log.Infof("Process %s daemon exit", sub.name)
        }()
    }
}

func killSub(sub *Process) {
    defer func() {
        recover()
    }()
    log.Infof("Send kill signal to process %s", sub.name)
    sub.sendSignal(syscall.SIGTERM)
}

func waitSub(sub *exec.Cmd) {
    defer func() {
        recover()
    }()
    _ = sub.Wait()
}

func exit() {
    exitFlag = true
    log.Warningf("Start exit process")
    for _, sub := range subProcess {
        if sub.process != nil {
            killSub(sub)
        }
    }
    log.Info("Wait all sub process exit")
    for _, sub := range subProcess {
        waitSub(sub.process)
    }
    log.Infof("Exit complete")
}
