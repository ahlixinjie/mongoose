package mongoose

import (
	"context"
	"errors"
	"github.com/ahlixinjie/mongoose/log"
	"github.com/ahlixinjie/mongoose/transport/common"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestLog(t *testing.T) {
	log.GetLogger().With(zap.String("hello", "world")).Info("can I")
	log.GetLogger().With(zap.Error(errors.New("custom error"))).Error("msg")

	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{
		common.HeaderRequestID: "123",
		log.Region:             "vn",
	}))
	log.GetLoggerWithCtx(ctx).Info("haha")
}
