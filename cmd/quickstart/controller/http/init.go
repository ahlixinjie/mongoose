package http

import (
	"context"
	api "github.com/ahlixinjie/fiction/api/http"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/dig"
	"google.golang.org/grpc"
)

type Impl struct {
	api.UnimplementedFictionServer
}

func (i *Impl) Provide() (constructor interface{}, _ []dig.ProvideOption) {
	type conf struct {
		dig.Out
		Handler       *runtime.ServeMux
		GatewayFunc   func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) (err error)
		FictionServer api.FictionServer `name:"http"`
	}

	constructor = func() conf {
		return conf{
			Handler:       runtime.NewServeMux(),
			GatewayFunc:   api.RegisterFictionHandlerFromEndpoint,
			FictionServer: i,
		}
	}
	return
}

func (i *Impl) Echo(ctx context.Context, request *api.StringMessage) (response *api.StringMessage, err error) {
	response = &api.StringMessage{Value: "from http " + request.GetValue()}
	return
}
