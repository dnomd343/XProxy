package main

import (
	"XProxy/next/assets"
	"XProxy/next/logger"
)

const gzSample = "/root/XProxy/LICENSE.gz"
const xzSample = "/root/XProxy/LICENSE.xz"
const bz2Sample = "/root/XProxy/LICENSE.bz2"

func main() {

	//raw, _ := os.ReadFile(gzSample)
	//raw, _ := os.ReadFile(bz2Sample)
	//raw, _ := os.ReadFile(xzSample)
	//Logger.Debugf("data len -> %d", len(raw))
	//ret, err := assets.Extract(raw)
	//if err != nil {
	//	Logger.Debugf("extract error -> %v", err)
	//}
	//Logger.Debugf("extract ok -> len = %d", len(ret))
	//os.WriteFile("demo.data", ret, 0777)

	//raw, _ := assets.Download("https://cdn.dnomd343.top/v2ray-rules-dat/geosite.dat", "")
	//raw, _ := assets.Download("https://cdn.dnomd343.top/v2ray-rules-dat/geosite.dat", "socks5://192.168.2.2:1084")
	raw, _ := assets.Download("https://cdn.dnomd343.top/v2ray-rules-dat/geosite.dat.xz", "")
	ret, _ := assets.Extract(raw)
	logger.Debugf("content size -> %d", len(ret))

}
