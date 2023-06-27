package mongoose

import (
	"github.com/ahlixinjie/mongoose/config"
	"github.com/ahlixinjie/mongoose/transport/grpc"
	"github.com/ahlixinjie/mongoose/transport/http"
	"go.uber.org/fx"
)

func Run(elements ...fx.Option) {
	elements = append(elements, fx.Provide(config.NewProvider),
		fx.Provide(grpc.NewGrpcServer), fx.Provide(http.NewHTTPServer), //grpc must init before http
	)
	app := fx.New(elements...)
	app.Run()
}
