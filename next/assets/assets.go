package assets

import (
	"XProxy/next/logger"
	"github.com/robfig/cron"
	urlpkg "net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

var buildinAssets = map[string]string{
	"geoip.dat":   "/geoip.dat.xz",
	"geosite.dat": "/geosite.dat.xz",
}

func LoadBuildin() {
	updateLocalAssets(buildinAssets, true)
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

func init() {
	updateChan := make(chan os.Signal, 1)
	go func() {
		for {
			<-updateChan
			logger.Debugf("Trigger assets update due to receiving SIGALRM")
			Update()
		}
	}()
	signal.Notify(updateChan, syscall.SIGALRM)
}

func mapClone(raw map[string]string) map[string]string {
	assets := make(map[string]string, len(raw))
	for file, url := range raw {
		assets[file] = strings.Clone(url)
	}
	return assets
}

func Update() {
	update.renew.Lock()
	proxy := update.proxy
	assets := mapClone(update.assets)
	update.renew.Unlock()

	if !update.running.TryLock() {
		logger.Infof("Another assets update is in progress, skip it")
		return
	}
	logger.Infof("Start remote assets update process")
	updateRemoteAssets(assets, proxy, false)
	update.running.Unlock()
}

func GetAssets() map[string]string {
	update.renew.Lock()
	assets := mapClone(update.assets)
	update.renew.Unlock()
	return assets
}

func SetAssets(assets map[string]string) {
	update.renew.Lock()
	update.assets = mapClone(assets)
	update.renew.Unlock()
}

func GetProxy() string {
	update.renew.Lock()
	proxy := update.proxy.String()
	update.renew.Unlock()
	return proxy
}

func SetProxy(proxy string) error {
	var proxyUrl *urlpkg.URL
	if proxy != "" {
		url, err := urlpkg.Parse(proxy)
		if err != nil {
			logger.Errorf("Invalid proxy url `%s` -> %v", proxy, err)
			return err
		}
		proxyUrl = url
	}
	update.renew.Lock()
	update.proxy = proxyUrl
	update.renew.Unlock()
	return nil
}

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
