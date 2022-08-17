package asset

import (
	"XProxy/cmd/common"
	log "github.com/sirupsen/logrus"
)

func extractGeoFile(archivePath string, geoFile string, targetDir string) {
	if common.IsFileExist(targetDir + "/" + geoFile) {
		log.Debugf("Asset %s exist -> skip extract", geoFile)
		return
	}
	log.Infof("Extract asset file -> %s", targetDir+"/"+geoFile)
	common.RunCommand("tar", "xvf", archivePath, "./"+geoFile, "-C", targetDir)
}

func LoadGeoIp(assetFile string, assetDir string) {
	common.CreateFolder(assetDir)
	extractGeoFile(assetFile, "geoip.dat", assetDir)
}

func LoadGeoSite(assetFile string, assetDir string) {
	common.CreateFolder(assetDir)
	extractGeoFile(assetFile, "geosite.dat", assetDir)
}
