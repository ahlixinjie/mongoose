package grpc

import (
	"fmt"
	"github.com/ahlixinjie/mongoose/log"
	"github.com/ahlixinjie/mongoose/transport"
	"github.com/ahlixinjie/mongoose/utils/parse"
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

const (
	confKeyRPC = "RPC"
)

type Service struct {
	port   int
	server *grpc.Server
}

func (s *Service) Provide() (constructor interface{}, opts []dig.ProvideOption) {
	type conf struct {
		dig.In
		Conf map[string]interface{} `name:"dig_conf"`
	}
	constructor = func(c conf) *grpc.Server {
		config := c.Conf
		v, ok := config[transport.ConfKeyPort]
		if !ok {
			return nil
		}

		vv, ok := v.(map[string]interface{})
		if !ok {
			return nil
		}
		portStr := vv[confKeyRPC].(string)
		s.port = parse.Port(portStr)
		if s.port == 0 {
			s.port = parse.Port(os.Getenv(portStr))
		}

		if s.port == 0 {
			panic("not set rpc port")
		}
		s.server = grpc.NewServer()
		return s.server
	}
	return
}

func (s *Service) Start() (err error) {
	if s.port == 0 {
		return
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
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
