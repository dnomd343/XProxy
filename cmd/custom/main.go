package custom

import (
    "XProxy/cmd/common"
    log "github.com/sirupsen/logrus"
)

type Config struct {
    Pre  []string `yaml:"pre" json:"pre"`
    Post []string `yaml:"post" json:"post"`
}

func RunPreScript(config *Config) {
    for _, script := range config.Pre {
        log.Infof("Run pre-script command -> %s", script)
        common.RunCommand("sh", "-c", script)
    }
}

func RunPostScript(config *Config) {
    for _, script := range config.Post {
        log.Infof("Run post-script command -> %s", script)
        common.RunCommand("sh", "-c", script)
    }
}
