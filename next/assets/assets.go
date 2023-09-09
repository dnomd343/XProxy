package assets

var buildinAssets = map[string]string{
	"geoip.dat":   "geoip.dat.xz",
	"geosite.dat": "geosite.dat.xz",
}

func Demo() {

	remoteAssets := map[string]string{
		"geoip.dat":   "https://cdn.dnomd343.top/v2ray-rules-dat/geoip.dat.xz",
		"geosite.dat": "https://cdn.dnomd343.top/v2ray-rules-dat/geosite.dat.xz",
	}

	//updateLocalAssets(buildinAssets, true)
	//updateLocalAssets(buildinAssets, false)

	updateRemoteAssets(remoteAssets, "", true)
	//updateRemoteAssets(remoteAssets, "", false)
	//updateRemoteAssets(remoteAssets, "socks5://192.168.2.2:1084", true)
	//updateRemoteAssets(remoteAssets, "socks5://192.168.2.2:1084", false)

	//time.Sleep(10 * time.Second)

	//updateRemoteAsset("geosite.dat", "http://cdn.dnomd343.top/v2ray-rules-dat/geosite.dat.xz", "")
	//updateRemoteAsset("geosite.dat", "http://cdn.dnomd343.top/v2ray-rules-dat/geosite.dat.xz", "socks5://192.168.2.2:1084")
	//updateLocalAsset("geosite.dat", "geosite.dat.xz")

}
