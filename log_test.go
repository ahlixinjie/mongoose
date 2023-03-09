package mongoose

import (
	"github.com/ahlixinjie/mongoose/log"
	"testing"
)

func TestLog(t *testing.T) {
	log.Info("hello")
	log.WithField("hello", "test").Info("11")
}
