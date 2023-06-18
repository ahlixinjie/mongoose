package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/ahlixinjie/mongoose/log"
	"github.com/ahlixinjie/mongoose/transport/common"
	transport_grpc "github.com/ahlixinjie/mongoose/transport/grpc"
	"github.com/ahlixinjie/mongoose/utils/parse"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/config"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
)

type Service struct {
	port    string
	handler http.Handler
}

var (
	port    int
	handler http.Handler
	server  *http.Server
)

type params struct {
	fx.In
	Lc          fx.Lifecycle
	Config      *config.YAML
	Handler     *runtime.ServeMux
	GatewayFunc func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) (err error)
}

func NewHTTPServer(p params) (err error) {
	portStr := p.Config.Get(common.ConfKeyPort + "." + common.ConfKeyHTTP).String()
	port = parse.Port(portStr)
	if port == 0 && len(portStr) != 0 {
		//try to get port from env
		port = parse.Port(os.Getenv(portStr))
	}
	if port == 0 {
		log.Info("won't start http service")
		return
	}
	if err = p.GatewayFunc(context.Background(), p.Handler, fmt.Sprintf(":%d", transport_grpc.GetPort()),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}); err != nil {
		return
	}

	handler = p.Handler
	p.Lc.Append(fx.StartHook(start))
	p.Lc.Append(fx.StopHook(stop))
	return
}

func start() {
	go func() {
		log.Info("start http service")
		server = &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: handler}
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Info("http service has been shutdown")
			return
		}
		if err != nil {
			panic(err)
		}
	}()
}

func stop(ctx context.Context) error {
	log.Info("trying to stop http service")
	return server.Shutdown(ctx)
}
