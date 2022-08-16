package main

import (
    log "github.com/sirupsen/logrus"
    "os"
    "os/exec"
    "syscall"
)

type Process struct {
    enable  bool
    caption string
    command []string
    process *exec.Cmd
}

func newProcess(command ...string) *Process {
    process := new(Process)
    process.enable = true
    process.command = command
    process.caption = command[0]
    log.Debugf("New process %s -> %v", process.caption, process.command)
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
        log.Errorf("Failed to start %s -> %v", p.caption, err)
    }
    log.Infof("Start process %s -> PID = %d", p.caption, p.process.Process.Pid)
}

func (p *Process) isProcessAlive() bool {
    return p.process.ProcessState == nil
}

func (p *Process) sendSignal(signal syscall.Signal) {
    err := p.process.Process.Signal(signal)
    if err != nil {
        log.Errorf("Send signal %v to process %s error -> %v", signal, p.caption, err)
    }
}

func (p *Process) waitProcess() {
    err := p.process.Wait()
    if err != nil {
        log.Warningf("Wait process %s -> %v", p.caption, err)
    }
}

func daemonSub(sub *Process) {
    for sub.isProcessAlive() {
        sub.waitProcess()
    }
    log.Warningf("Catch process %s exit", sub.caption)
    if sub.enable {
        sub.startProcess(true, true)
        log.Infof("Process %s restart success", xray.caption)
        daemonSub(sub)
    }
}

func daemon(sub *Process) {
    go func() {
        daemonSub(sub)
    }()
}
