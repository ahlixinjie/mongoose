package mongoose

import (
	"errors"
	"github.com/ahlixinjie/mongoose/log"
	"go.uber.org/zap"
	"testing"
)

func TestLog(t *testing.T) {
	log.GetLogger().With(zap.String("hello", "world")).Info("can I")
	log.GetLogger().With(zap.Error(errors.New("custom error"))).Error("msg")
}
