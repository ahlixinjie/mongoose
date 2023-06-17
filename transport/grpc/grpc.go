package grpc

import (
	"fmt"
	"github.com/ahlixinjie/mongoose/log"
	"github.com/ahlixinjie/mongoose/transport/common"
	"github.com/ahlixinjie/mongoose/utils/parse"
	"go.uber.org/config"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

var (
	port   int
	server *grpc.Server
)

func NewGrpcServer(lc fx.Lifecycle, config *config.YAML) *grpc.Server {
	portStr := config.Get(common.ConfKeyPort + "." + common.ConfKeyRPC).String()
	port = parse.Port(portStr)
	if port == 0 && len(portStr) != 0 {
		//try to get port from env
		port = parse.Port(os.Getenv(portStr))
	}
	if port == 0 {
		log.Info("won't start grpc service")
		return nil
	}
	server = grpc.NewServer()

	lc.Append(fx.StartHook(start))
	lc.Append(fx.StopHook(stop))
	return server
}

func start() (err error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	reflection.Register(server)
	go func() {
		log.Info("start grpc service")
		err := server.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()

	return nil
}

func stop() {
	log.Info("trying to stop grpc service")
	server.GracefulStop()
}

func GetPort() int {
	return port
}
