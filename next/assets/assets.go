package assets

import "XProxy/next/logger"

func Demo() {
	raw, _ := download("https://cdn.dnomd343.top/v2ray-rules-dat/geosite.dat", "")
	//raw, _ := download("https://cdn.dnomd343.top/v2ray-rules-dat/geosite.dat", "socks5://192.168.2.2:1084")
	//raw, _ := download("https://cdn.dnomd343.top/v2ray-rules-dat/geosite.dat.xz", "")
	ret, _ := tryExtract(raw)
	logger.Debugf("content size -> %d", len(ret))
}
