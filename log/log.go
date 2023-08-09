package log

import (
	"context"
	"github.com/ahlixinjie/go-utils/collections/list"
	"github.com/ahlixinjie/mongoose/transport/common"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/metadata"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
)

const (
	Region = "fiction-region"
)

var logger *zap.Logger

func init() {
	writeOnFileCore := zapcore.NewCore(getEncoder(), getLogWriteSyncer("log/all.log"), getLogLevel())

	tee := zapcore.NewTee(writeOnFileCore)

	logger = zap.New(tee, zap.WithCaller(true))
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

func GetLoggerWithCtx(ctx context.Context) *zap.Logger {
	var l = logger
	if region := list.GetFirstElement(metadata.ValueFromIncomingContext(ctx, Region)); len(region) != 0 {
		l = l.With(zap.String(Region, region))
	}
	if requestID := list.GetFirstElement(metadata.ValueFromIncomingContext(ctx, strings.ToLower(common.HeaderRequestID))); len(requestID) != 0 {
		l = l.With(zap.String(common.HeaderRequestID, requestID))
	}
	return l
}
