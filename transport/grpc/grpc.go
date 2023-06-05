package grpc

import (
	"fmt"
	"github.com/ahlixinjie/mongoose/log"
	"github.com/ahlixinjie/mongoose/transport/common"
	"github.com/ahlixinjie/mongoose/utils/parse"
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

type Service struct {
	port   string
	server *grpc.Server
}

func (s *Service) Provide() (constructor interface{}, _ []dig.ProvideOption) {
	type conf struct {
		dig.In
		Conf map[string]interface{} `name:"dig_conf"`
	}
	constructor = func(c conf) *grpc.Server {
		config := c.Conf
		v, ok := config[common.ConfKeyPort]
		if !ok {
			return nil
		}

		vv, ok := v.(map[string]interface{})
		if !ok {
			return nil
		}
		s.port = vv[common.ConfKeyRPC].(string)
		if parse.Port(s.port) == 0 {
			s.port = os.Getenv(s.port)
		}

		if len(s.port) == 0 {
			return nil
		}
		s.server = grpc.NewServer()
		return s.server
	}
	return
}

func (s *Service) Start() (err error) {
	if len(s.port) == 0 {
		fmt.Println("not set rpc service")
		return
	}
	lis, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}
	reflection.Register(s.server)
	go func() {
		log.Info("start grpc")
		err := s.server.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()

	return nil
}
