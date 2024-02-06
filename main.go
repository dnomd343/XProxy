package main

import "XProxy/assets"

func main() {
	remoteAssets := map[string]string{
		"geoip.dat":   "https://cdn.dnomd343.top/v2ray-rules-dat/geoip.dat.xz",
		"geosite.dat": "https://cdn.dnomd343.top/v2ray-rules-dat/geosite.dat.xz",
	}
	assets.SetAssets(remoteAssets)

	//assets.Update()

	select {}

	//logger.Warnf("cron -> `%s`", assets.GetCron())
	//assets.SetCron("@every 1s")
	//logger.Warnf("cron -> `%s`", assets.GetCron())
	//
	//time.Sleep(5 * time.Second)
	//assets.SetCron("@every 2s")
	//logger.Warnf("cron -> `%s`", assets.GetCron())
	//
	//time.Sleep(8 * time.Second)
	//assets.SetCron("")
	//logger.Warnf("cron -> `%s`", assets.GetCron())
	//select {}

	//assets.LoadBuildin()

	//assets.Update(false)
	//assets.SetUpdateConfig(assets.UpdateSettings{
	//	cron: "",
	//	ass
	//})
}
