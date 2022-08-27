package asset

import (
    "XProxy/cmd/common"
    "github.com/robfig/cron"
    log "github.com/sirupsen/logrus"
    "path"
)

type Config struct {
    Proxy string            `yaml:"proxy" json:"proxy"`
    Cron  string            `yaml:"cron" json:"cron"`
    Url   map[string]string `yaml:"url" json:"url"`
}

func updateAsset(urls map[string]string, assetDir string, updateProxy string) { // download new assets
    defer func() {
        if err := recover(); err != nil {
            log.Errorf("Update failed -> %v", err)
        }
    }()
    if len(urls) != 0 {
        log.Info("Start update assets")
        for file, url := range urls {
            common.DownloadFile(url, path.Join(assetDir, file), updateProxy) // maybe override old asset
        }
    }
}

func AutoUpdate(update *Config, assetDir string) { // set cron task for auto update
    if update.Cron != "" {
        autoUpdate := cron.New()
        _ = autoUpdate.AddFunc(update.Cron, func() { // cron function
            updateAsset(update.Url, assetDir, update.Proxy)
        })
        autoUpdate.Start()
    }
}
