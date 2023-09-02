package logger

import "go.uber.org/zap"

var Logger *zap.SugaredLogger

func init() {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic("log utils init failed")
	}

	// TODO: more zap logger configure
	// TODO: reserve raw logger handle

	Logger = log.Sugar()
}
