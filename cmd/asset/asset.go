package asset

import (
    "XProxy/cmd/common"
    log "github.com/sirupsen/logrus"
)

func extractFile(archive string, geoFile string, targetDir string) {
    if common.IsFileExist(targetDir + "/" + geoFile) {
        log.Debugf("Asset %s exist -> skip extract", geoFile)
        return
    }
    log.Infof("Extract asset file -> %s", targetDir+"/"+geoFile)
    common.RunCommand("tar", "xvf", archive, "./"+geoFile, "-C", targetDir)
}

func Load(assetFile string, assetDir string) {
    common.CreateFolder(assetDir)
    extractFile(assetFile, "geoip.dat", assetDir)
    extractFile(assetFile, "geosite.dat", assetDir)
}
