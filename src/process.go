package main

import (
    log "github.com/sirupsen/logrus"
    "os"
    "os/exec"
    "syscall"
    "time"
)

var subProcess []*Process

type Process struct {
    enable  bool
    name    string
    command []string
    process *exec.Cmd
    exit    bool
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
    p.enable = true
    err := p.process.Start()
    if err != nil {
        log.Errorf("Failed to start %s -> %v", p.name, err)
    }
    log.Infof("Start process %s -> PID = %d", p.name, p.process.Process.Pid)
}

func (p *Process) sendSignal(signal syscall.Signal) {
    //defer func() {
    //    _ = recover()
    //}()
    if p.process != nil && p.process.ProcessState == nil {
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
    if !sub.enable {
        log.Debugf("Process %s disabled -> stop daemon", sub.name)
        return
    }
    sub.startProcess(true, true)
    log.Infof("Process %s restart success", sub.name)
    daemonSub(sub)
}

func daemon(sub *Process) {
    if !sub.enable {
        log.Infof("Process %s disabled -> skip daemon", sub.name)
        sub.exit = true
        return
    }
    log.Infof("Start daemon of process %s", sub.name)
    go func() {
        daemonSub(sub)
        log.Infof("Process %s daemon exit", sub.name)
        sub.exit = true
    }()
}

func exit() {
    log.Warningf("Start exit process")
    for _, sub := range subProcess {
        sub.enable = false
        if sub.process != nil {
            sub.sendSignal(syscall.SIGTERM)
            log.Infof("Send kill signal to process %s", sub.name)
        }
    }
    var allExit bool
    log.Info("Wait all sub process exit")
    for !allExit {
        time.Sleep(10 * time.Millisecond) // delay 10ms
        allExit = true
        for _, sub := range subProcess {
            allExit = allExit && sub.exit //(sub.process.ProcessState != nil)
        }
    }
    log.Infof("Exit complete")
}
