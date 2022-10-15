package asset

import (
    "XProxy/cmd/common"
    "github.com/robfig/cron"
    log "github.com/sirupsen/logrus"
    "path"
)

type Config struct {
    Disable bool `yaml:"disable" json:"disable" toml:"disable"`
    Update  struct {
        Proxy string            `yaml:"proxy" json:"proxy" toml:"proxy"`
        Cron  string            `yaml:"cron" json:"cron" toml:"cron"`
        Url   map[string]string `yaml:"url" json:"url" toml:"url"`
    }
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
            if !common.DownloadFile(url, path.Join(assetDir, file), updateProxy) { // maybe override old asset
                log.Infof("Try to download asset `%s` again", file)
                common.DownloadFile(url, path.Join(assetDir, file), updateProxy) // download retry
            }
        }
        log.Infof("Assets update complete")
    }
}

func AutoUpdate(config *Config, assetDir string) { // set cron task for auto update
    if config.Update.Cron != "" {
        autoUpdate := cron.New()
        _ = autoUpdate.AddFunc(config.Update.Cron, func() { // cron function
            updateAsset(config.Update.Url, assetDir, config.Update.Proxy)
        })
        autoUpdate.Start()
    }
}
