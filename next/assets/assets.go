package assets

import (
	"XProxy/next/logger"
	urlpkg "net/url"
	"strings"
	"sync"
)

var buildinAssets = map[string]string{
	"geoip.dat":   "/geoip.dat.xz",
	"geosite.dat": "/geosite.dat.xz",
}

type UpdateSettings struct {
	cron   string
	mutex  sync.Mutex
	proxy  *urlpkg.URL
	assets map[string]string
}

var update UpdateSettings

func assetsClone(raw map[string]string) map[string]string {
	assets := make(map[string]string, len(raw))
	for file, url := range raw {
		assets[file] = strings.Clone(url)
	}
	return assets
}

func SetCron(cron string) error {
	// TODO: setting up crond service
	return nil
}

func GetProxy() string {
	update.mutex.Lock()
	proxy := update.proxy.String()
	update.mutex.Unlock()
	return proxy
}

func SetProxy(proxy string) error {
	var proxyUrl *urlpkg.URL // clear proxy by empty string
	if proxy != "" {
		url, err := urlpkg.Parse(proxy)
		if err != nil {
			logger.Errorf("Invalid proxy url `%s` -> %v", proxy, err)
			return err
		}
		proxyUrl = url
	}
	update.mutex.Lock()
	update.proxy = proxyUrl
	update.mutex.Unlock()
	return nil
}

func SetAssets(assets map[string]string) {
	update.mutex.Lock()
	update.assets = assetsClone(assets)
	update.mutex.Unlock()
}

func GetAssets() map[string]string {
	update.mutex.Lock()
	assets := assetsClone(update.assets)
	update.mutex.Unlock()
	return assets
}

func LoadBuildin() {
	updateLocalAssets(buildinAssets, true)
}
