package common

import (
	log "github.com/sirupsen/logrus"
	"os/exec"
	"syscall"
)

func RunCommand(command ...string) (int, string) {
	log.Debugf("Running system command -> %v", command)
	process := exec.Command(command[0], command[1:]...)
	output, _ := process.CombinedOutput()
	log.Debugf("Command %v -> \n%s", command, string(output))
	code := process.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	if code != 0 {
		log.Warningf("Command %v return code %d", command, code)
	}
	return code, string(output)
}
