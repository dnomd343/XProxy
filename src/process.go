package main

import (
    log "github.com/sirupsen/logrus"
    "os"
    "os/exec"
    "syscall"
)

type Process struct {
    command []string
    process *exec.Cmd
}

func newProcess(command ...string) *Process {
    process := new(Process)
    process.command = command
    log.Debugf("New process -> %v", command)
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
        log.Errorf("Failed to start %v -> %v", p.command, err)
    }
    log.Infof("Start process %v -> PID = %d", p.command, p.process.Process.Pid)
}

func (p *Process) isProcessAlive() bool {
    return p.process.ProcessState == nil
}

func (p *Process) sendSignal(signal syscall.Signal) {
    err := p.process.Process.Signal(signal)
    if err != nil {
        log.Errorf("Send signal %v error -> %v", signal, p.command)
    }
}

func (p *Process) waitProcess() {
    err := p.process.Wait()
    if err != nil {
        log.Errorf("Wait process error -> %v", p.command)
    }
}
