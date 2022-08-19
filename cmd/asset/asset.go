package asset

import (
    "XProxy/cmd/common"
    log "github.com/sirupsen/logrus"
    "path"
)

func extractFile(archive string, geoFile string, targetDir string) { // extract `.dat` file into targetDir
    filePath := path.Join(targetDir, geoFile)
    if common.IsFileExist(filePath) {
        log.Debugf("Asset %s exist -> skip extract", geoFile)
        return
    }
    log.Infof("Extract asset file -> %s", filePath)
    common.RunCommand("tar", "xvf", archive, "./"+geoFile, "-C", targetDir)
}

func Load(assetFile string, assetDir string) {
    common.CreateFolder(assetDir)
    extractFile(assetFile, "geoip.dat", assetDir)
    extractFile(assetFile, "geosite.dat", assetDir)
}
