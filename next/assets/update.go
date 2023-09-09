package assets

import (
	"XProxy/next/logger"
	"os"
	"sync"
	"time"
)

// saveAsset is used to write to a local file and specify its last modify
// time. If the file exists, it will be replaced.
func saveAsset(file string, content []byte, date *time.Time) error {
	// TODO: create folder if not exist
	fp, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		logger.Errorf("Failed to open file `%s` -> %v", file, err)
		return err
	}
	defer fp.Close()

	_, err = fp.Write(content)
	if err != nil {
		logger.Errorf("Failed to save file `%s` -> %v", file, err)
		return err
	}

	if date != nil { // change local file last modify time
		if err := os.Chtimes(file, *date, *date); err != nil {
			logger.Warnf("Failed to change asset modification time")
		} else {
			logger.Debugf("Change `%s` modification time to `%v`", file, date)
		}
	}
	return nil
}

// updateRemoteAsset will download remote asset via the optional proxy and
// save them locally. Local files will be overwritten if they exist.
func updateRemoteAsset(file string, url string, proxy string) error {
	logger.Debugf("Start downloading remote asset `%s` to `%s`", url, file)
	asset, date, err := downloadAsset(url, proxy)
	if err != nil {
		logger.Errorf("Failed to download remote asset `%s`", url)
		return err
	}
	if asset, err = tryExtract(asset); err != nil {
		return err
	}
	if err := saveAsset(file, asset, date); err != nil {
		return err
	}
	logger.Infof("Successfully obtained remote asset `%s`", file)
	return nil
}

// updateLocalAsset will extract local asset and save them locally. If the
// local file already exists, it will be skipped.
func updateLocalAsset(file string, src string) error {
	_, err := os.Stat(file)
	if err == nil {
		logger.Debugf("Local asset `%s` already exist", file)
		return nil // skip local asset extract
	}

	var date = time.Now()
	if stat, err := os.Stat(src); err == nil {
		date = stat.ModTime() // using last modify time of src file
	}

	logger.Debugf("Start extracting local asset `%s`", file)
	asset, err := os.ReadFile(src)
	if err != nil {
		logger.Errorf("Failed to read local asset -> %v", err)
		return err
	}
	if asset, err = tryExtract(asset); err != nil {
		return err
	}
	if err := saveAsset(file, asset, &date); err != nil {
		return err
	}
	logger.Infof("Successfully extracted local asset `%s`", file)
	return nil
}

// updateMultiple performs multiple update tasks, optionally in parallel.
func updateMultiple(updates []func(), parallel bool) {
	if !parallel {
		for _, update := range updates {
			update() // sequential update
		}
		return
	}
	var wg sync.WaitGroup
	for _, update := range updates {
		update_ := update
		wg.Add(1)
		go func() { // parallel update
			update_()
			wg.Done()
		}()
	}
	wg.Wait() // wait all update tasks done
}

// updateRemoteAssets will update multiple remote assets in batches.
func updateRemoteAssets(assets map[string]string, proxy string, parallel bool) {
	var tasks []func()
	for file, url := range assets {
		file_, url_ := file, url
		tasks = append(tasks, func() {
			_ = updateRemoteAsset(file_, url_, proxy)
		})
	}
	updateMultiple(tasks, parallel)
}

// updateLocalAssets will extract multiple local assets in batches.
func updateLocalAssets(assets map[string]string, parallel bool) {
	var tasks []func()
	for file, src := range assets {
		file_, src_ := file, src
		tasks = append(tasks, func() {
			_ = updateLocalAsset(file_, src_)
		})
	}
	updateMultiple(tasks, parallel)
}
