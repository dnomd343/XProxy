package assets

import "XProxy/next/logger"

func Demo() {
	//raw, t, _ := download("https://github.com/Loyalsoldier/v2ray-rules-dat/releases/download/202309082208/geosite.dat", "")
	//raw, t, _ := download("https://cdn.dnomd343.top/v2ray-rules-dat/geosite.dat", "")
	//raw, t, _ := download("https://cdn.dnomd343.top/v2ray-rules-dat/geosite.dat", "socks5://192.168.2.2:1084")
	raw, t, _ := download("https://cdn.dnomd343.top/v2ray-rules-dat/geosite.dat.xz", "")
	logger.Infof("%v", t)
	ret, _ := tryExtract(raw)
	logger.Debugf("content size -> %d", len(ret))
}
