package http

import (
	"context"
	"fmt"
	"github.com/ahlixinjie/mongoose/log"
	"github.com/ahlixinjie/mongoose/transport/common"
	"github.com/ahlixinjie/mongoose/utils/parse"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
)

type Service struct {
	port    string
	handler http.Handler
}

func (s *Service) Invoke() (function interface{}, _ []dig.InvokeOption) {
	type conf struct {
		dig.In
		Conf        map[string]interface{} `name:"dig_conf"`
		Handler     *runtime.ServeMux
		GatewayFunc func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) (err error)
	}

	function = func(c conf) {
		config := c.Conf
		v, ok := config[common.ConfKeyPort]
		if !ok {
			return
		}
		vv, ok := v.(map[string]interface{})
		if !ok {
			return
		}

		s.port = vv[common.ConfKeyHTTP].(string)
		if parse.Port(s.port) == 0 {
			s.port = os.Getenv(s.port)
		}

		if len(s.port) == 0 {
			return
		}

		if err := c.GatewayFunc(context.Background(), c.Handler, vv[common.ConfKeyRPC].(string),
			[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}); err != nil {
			panic(err)
		}

		s.handler = c.Handler
	}
	return
}

func (s *Service) Start() (err error) {
	if len(s.port) == 0 {
		fmt.Println("not set HTTP service")
		return
	}

	go func() {
		log.Info("start http")
		err := http.ListenAndServe(s.port, s.handler)
		if err != nil {
			panic(err)
		}
	}()

	return nil
}
