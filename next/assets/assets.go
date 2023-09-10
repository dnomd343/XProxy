package assets

import (
	"XProxy/next/logger"
	"github.com/robfig/cron"
	urlpkg "net/url"
	"strings"
	"sync"
)

var buildinAssets = map[string]string{
	"geoip.dat":   "/geoip.dat.xz",
	"geosite.dat": "/geosite.dat.xz",
}

type updateConfig struct {
	spec    string
	cron    *cron.Cron
	renew   sync.Mutex
	running sync.Mutex
	proxy   *urlpkg.URL
	assets  map[string]string
}

var update updateConfig

//func assetsClone(raw map[string]string) map[string]string {
//	assets := make(map[string]string, len(raw))
//	for file, url := range raw {
//		assets[file] = strings.Clone(url)
//	}
//	return assets
//}

// GetCron is used to obtain cron service specification.
func GetCron() string {
	update.renew.Lock()
	spec := strings.Clone(update.spec)
	update.renew.Unlock()
	return spec
}

// SetCron is used to update cron service specification.
func SetCron(spec string) error {
	if spec == update.spec {
		return nil // cron spec without renew
	}

	var cs *cron.Cron
	if spec != "" { // update cron service
		cs = cron.New()
		err := cs.AddFunc(spec, func() {
			var entry *cron.Entry
			if entries := update.cron.Entries(); len(entries) != 0 && entries[0] != nil {
				entry = entries[0]
			}
			logger.Debugf("hello from cron")
			if entry != nil {
				logger.Debugf("Assets cron service next trigger -> `%s`", entry.Next)
			}
		})
		if err != nil {
			logger.Errorf("Invalid cron spec `%s` -> %v", spec, err)
			return err
		}
		cs.Start()
	}

	update.renew.Lock()
	if update.cron != nil {
		update.cron.Stop() // stop old cron service
	}
	update.cron = cs
	update.spec = spec
	if cs == nil {
		logger.Infof("Assets cron service has been terminated")
	} else {
		logger.Infof("Assets cron service has been updated -> `%s`", spec)
	}
	update.renew.Unlock()
	return nil
}

//func GetProxy() string {
//	update.mutex.Lock()
//	proxy := update.proxy.String()
//	update.mutex.Unlock()
//	return proxy
//}
//
//func SetProxy(proxy string) error {
//	var proxyUrl *urlpkg.URL // clear proxy by empty string
//	if proxy != "" {
//		url, err := urlpkg.Parse(proxy)
//		if err != nil {
//			logger.Errorf("Invalid proxy url `%s` -> %v", proxy, err)
//			return err
//		}
//		proxyUrl = url
//	}
//	update.mutex.Lock()
//	update.proxy = proxyUrl
//	update.mutex.Unlock()
//	return nil
//}
//
//func SetAssets(assets map[string]string) {
//	update.mutex.Lock()
//	update.assets = assetsClone(assets)
//	update.mutex.Unlock()
//}
//
//func GetAssets() map[string]string {
//	update.mutex.Lock()
//	assets := assetsClone(update.assets)
//	update.mutex.Unlock()
//	return assets
//}

func LoadBuildin() {
	updateLocalAssets(buildinAssets, true)
}
