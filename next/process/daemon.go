package process

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

const (
	Ready   = iota // process ready to start
	Running        // process is running for now
	Death          // process exited and waiting for restart
	Stopped        // process has been stopped
)

type Process struct {
	name   string
	args   []string
	envs   []string
	mutex  sync.Mutex
	proc   *exec.Cmd
	status int
}

func NewDaemon(name string, command []string, env map[string]string) *Process {
	fmt.Printf("Create new daemon task %s -> %v %v\n", name, command, env)
	var envs []string
	for k, v := range env {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}
	return &Process{
		name:   name,
		envs:   envs,
		args:   command,
		status: Ready,
	}
}

func (p *Process) boot() bool {
	p.proc = exec.Command(p.args[0], p.args[1:]...)
	p.proc.Stdout = os.Stdout
	p.proc.Stderr = os.Stderr
	if len(p.envs) != 0 { // with custom env variables
		p.proc.Env = p.envs
	}
	err := p.proc.Start()
	if err != nil {
		fmt.Printf("Failed to boot process %s -> %v\n", p.name, err)
		p.status = Stopped
		return false
	}
	fmt.Printf("Successfully executed process %s -> PID = %d\n", p.name, p.proc.Process.Pid)
	p.status = Running
	return true
}

func (p *Process) wait() {
	for p.proc.ProcessState == nil { // wait until exit
		p.proc.Wait()
	}
	fmt.Printf("Catch process %s exit", p.name)
}

func (p *Process) Start() {
	if !p.mutex.TryLock() {
		fmt.Printf("already running")
		return
	}

	// TODO: daemon logical
	p.boot()
	p.wait()

}

func (p *Process) Status() int {
	return p.status
}

func (p *Process) Stop() {
	// TODO: stop daemon and kill process
}
