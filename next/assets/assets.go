package assets

import (
	"XProxy/next/logger"
	"os"
)

func updateRemoteAsset(file string, url string, proxy string) error {
	logger.Debugf("Start downloading remote asset `%s` to `%s`", url, file)
	asset, date, err := download(url, proxy)
	if err != nil {
		logger.Errorf("Failed to download remote asset `%s`", url)
		return err
	}
	asset, err = tryExtract(asset)

	fp, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		logger.Errorf("Failed to open file `%s` -> %v", file, err)
		return err
	}
	defer fp.Close()
	_, err = fp.Write(asset)
	if err != nil {
		logger.Errorf("Failed to save file `%s` -> %v", file, err)
		return err
	}

	if err := os.Chtimes(file, *date, *date); err != nil {
		logger.Warnf("Failed to change asset modification time")
	} else {
		logger.Debugf("Change `%s` modification time to `%v`", file, date)
	}
	logger.Infof("Successfully obtained remote asset `%s`", file)
	return nil
}

func Demo() {
	//raw, t, _ := download("https://github.com/Loyalsoldier/v2ray-rules-dat/releases/download/202309082208/geosite.dat", "")
	//raw, t, _ := download("https://cdn.dnomd343.top/v2ray-rules-dat/geosite.dat", "")
	//raw, t, _ := download("https://cdn.dnomd343.top/v2ray-rules-dat/geosite.dat", "socks5://192.168.2.2:1084")
	//raw, t, _ := download("https://cdn.dnomd343.top/v2ray-rules-dat/geosite.dat.xz", "")
	//logger.Infof("%v", t)
	//ret, _ := tryExtract(raw)
	//logger.Debugf("content size -> %d", len(ret))

	updateRemoteAsset("geosite.dat", "http://cdn.dnomd343.top/v2ray-rules-dat/geosite.dat.xz", "")

}
