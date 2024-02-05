package main

import (
	"XProxy/logger"
	"os"
)

func main() {
	//logger.SetLevel(zap.DebugLevel)

	//logger.Debugf("here is %s level", "debug")
	//logger.Infof("here is %s level", "info")
	//logger.Warnf("here is %s level", "warn")
	//logger.Errorf("here is %s level", "error")
	//logger.Panicf("here is %s level", "panic")

	fp1, _ := os.Create("demo_1.log")
	fp2, _ := os.Create("demo_2.log")
	fp3, _ := os.Create("demo_3.log")

	logger.Debugf("output msg 1 at debug")
	logger.Infof("output msg 1 at info")
	logger.Warnf("output msg 1 at warn")
	logger.AddWriters(false, fp1, fp2)
	logger.SetLevel(logger.InfoLevel)
	logger.Debugf("output msg 2 at debug")
	logger.Infof("output msg 2 at info")
	logger.Warnf("output msg 2 at warn")
	logger.SetLevel(logger.WarnLevel)
	logger.AddWriters(true, fp3)
	logger.Debugf("output msg 3 at debug")
	logger.Infof("output msg 3 at info")
	logger.Warnf("output msg 3 at warn")
	logger.SetLevel(logger.DebugLevel)
	logger.Debugf("output msg 4 at debug")
	logger.Infof("output msg 4 at info")
	logger.Warnf("output msg 4 at warn")
}
