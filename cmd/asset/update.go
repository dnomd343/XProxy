package asset

import (
    "XProxy/cmd/common"
    "github.com/robfig/cron"
    log "github.com/sirupsen/logrus"
)

type Config struct {
    Cron string            `yaml:"cron" json:"cron"`
    Url  map[string]string `yaml:"url" json:"url"`
}

func updateAsset(urls map[string]string, assetDir string) {
    if len(urls) != 0 {
        log.Info("Start update assets")
        for file, url := range urls {
            common.DownloadFile(url, assetDir+"/"+file)
        }
    }
}

func AutoUpdate(update *Config, assetDir string) {
    autoUpdate := cron.New()
    _ = autoUpdate.AddFunc(update.Cron, func() {
        updateAsset(update.Url, assetDir)
    })
    autoUpdate.Start()
}
