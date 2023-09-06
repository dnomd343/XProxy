package main

import "XProxy/next/logger"

func main() {

	logger.Debugf("here is %s level", "debug")
	logger.Infof("here is %s level", "info")
	logger.Warnf("here is %s level", "warn")
	logger.Errorf("here is %s level", "error")
	//logger.Panicf("here is %s level", "panic")

}
