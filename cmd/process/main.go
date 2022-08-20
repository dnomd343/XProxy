package process

import (
    log "github.com/sirupsen/logrus"
    "os"
    "os/exec"
    "syscall"
)

type Process struct {
    name    string
    env     []string
    command []string
    process *exec.Cmd
}

func New(command ...string) *Process {
    process := new(Process)
    process.name = command[0]
    process.command = command
    log.Debugf("New process %s -> %v", process.name, process.command)
    return process
}

func (p *Process) Run(isOutput bool, env []string) {
    p.process = exec.Command(p.command[0], p.command[1:]...)
    if isOutput {
        p.process.Stdout = os.Stdout
        p.process.Stderr = os.Stderr
    }
    p.env = env
    if len(p.env) != 0 {
        p.process.Env = p.env
        log.Infof("Process %s with env -> %v", p.name, p.env)
    }
    err := p.process.Start()
    if err != nil {
        log.Errorf("Failed to start %s -> %v", p.name, err)
    }
    log.Infof("Start process %s -> PID = %d", p.name, p.process.Process.Pid)
}

func (p *Process) Signal(signal syscall.Signal) {
    if p.process != nil {
        log.Debugf("Send signal %v to %s", signal, p.name)
        _ = p.process.Process.Signal(signal)
    }
}

func (p *Process) Wait() {
    if p.process != nil {
        err := p.process.Wait()
        if err != nil {
            log.Warningf("Wait process %s -> %v", p.name, err)
        }
    }
}
