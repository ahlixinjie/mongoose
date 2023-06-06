package grpc

import (
	api "github.com/ahlixinjie/fiction/api/http"
	"go.uber.org/dig"
	"google.golang.org/grpc"
)

type Impl struct {
	api.FictionServer
}

func (i *Impl) Invoke() (function interface{}, opts []dig.InvokeOption) {
	type conf struct {
		dig.In
		Server        *grpc.Server
		FictionServer api.FictionServer `name:"http"`
	}
	function = func(conf conf) {
		i.FictionServer = conf.FictionServer
		api.RegisterFictionServer(conf.Server, i)
	}
	return
}
