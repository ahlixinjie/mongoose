package log

import "go.uber.org/zap"

var logger *zap.Logger

func init() {
	l, err := zap.NewDevelopment(zap.AddCaller())
	if err != nil {
		panic(err)
	}
	logger = l
}

func GetLogger() *zap.Logger {
	return logger
}
