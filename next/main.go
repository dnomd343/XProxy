package main

import (
	. "XProxy/next/logger"
)

func main() {
	Logger.Debugf("here is %s level", "debug")
	Logger.Infof("here is %s level", "info")
	Logger.Warnf("here is %s level", "warn")
	Logger.Errorf("here is %s level", "error")
	//Logger.Panicf("here is %s level", "panic")
}
