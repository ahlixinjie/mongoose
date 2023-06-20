package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var logger *zap.Logger

func init() {
	writeOnFileCore := zapcore.NewCore(getEncoder(), getLogWriteSyncer("log/all.log"), getLogLevel())

	tee := zapcore.NewTee(writeOnFileCore)

	logger = zap.New(tee)
}

func getEncoder() zapcore.Encoder {
	config := zap.NewDevelopmentEncoderConfig()
	config.FunctionKey = "F"
	config.ConsoleSeparator = "|"
	return zapcore.NewConsoleEncoder(config)
}

func getLogLevel() zapcore.Level {
	logLevel, err := zapcore.ParseLevel(os.Getenv("log-level"))
	if err != nil {
		logLevel = zapcore.DebugLevel
	}
	return logLevel
}

func getLogWriteSyncer(filePath string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:  filePath,
		MaxSize:   5,
		LocalTime: true,
		Compress:  false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func GetLogger() *zap.Logger {
	return logger
}
