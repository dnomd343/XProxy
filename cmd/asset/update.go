package asset

import (
	"XProxy/cmd/common"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

func updateAssets(urls map[string]string, assetDir string) {
	if len(urls) != 0 {
		log.Info("Start update assets")
		for file, url := range urls {
			common.DownloadFile(url, assetDir+"/"+file)
		}
	}
}

func AutoUpdate(updateCron string, updateUrls map[string]string, assetDir string) {
	autoUpdate := cron.New()
	_ = autoUpdate.AddFunc(updateCron, func() {
		updateAssets(updateUrls, assetDir)
	})
	autoUpdate.Start()
}
