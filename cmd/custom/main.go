package custom

import (
    log "github.com/sirupsen/logrus"
    "os"
    "os/exec"
)

type Config struct {
    Pre  []string `yaml:"pre" json:"pre"`
    Post []string `yaml:"post" json:"post"`
}

func runScript(command string) {
    log.Debugf("Run script -> %s", command)
    cmd := exec.Command("sh", "-c", command)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Start()
    if err != nil {
        log.Warningf("Script `%s` working error", command)
    } else {
        _ = cmd.Wait()
    }
}

func RunPreScript(config *Config) {
    for _, script := range config.Pre {
        log.Infof("Run pre-script command -> %s", script)
        runScript(script)
    }
}

func RunPostScript(config *Config) {
    for _, script := range config.Post {
        log.Infof("Run post-script command -> %s", script)
        runScript(script)
    }
}
