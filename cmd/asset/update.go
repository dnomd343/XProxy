package asset

import (
	"XProxy/cmd/common"
	log "github.com/sirupsen/logrus"
)

func UpdateAssets(urls map[string]string, assetDir string) {
	if len(urls) != 0 {
		log.Info("Start update assets")
		for file, url := range urls {
			common.DownloadFile(url, assetDir+"/"+file)
		}
	}
}
