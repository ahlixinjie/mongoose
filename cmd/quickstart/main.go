package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

//todo gitignore api/http/*.pb.* .idea/
//todo add make file
//todo 改写grpc && http controller

var dirNameFlag = flag.String("d", "demo", "input the directory name here")
var moduleNamFlag = flag.String("m", "github.com/ahlixinjie/demo", "input module name here")

func main() {
	flag.Parse()

	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("execute error: %v\n", e)
		}
	}()

	_, err := exec.Command("mkdir", *dirNameFlag).Output()
	if err != nil {
		panic(err)
	}

	if _, err = setCommandDir(exec.Command("go", "mod", "init", *moduleNamFlag), *dirNameFlag).Output(); err != nil {
		panic(err)
	}

	if _, err = setCommandDir(exec.Command("mkdir", "-p",
		"cmd/run", "api/http", "config",
		"internal/service", "internal/controller/grpc", "internal/controller/http"),
		*dirNameFlag).Output(); err != nil {
		panic(err)
	}

	writeProto()
	writeBufConfig()
	writeServiceConfig()
	writeGrpcController()
	writeHttpController()
	writeMain()
	if _, err = setCommandDir(exec.Command("go", "mod", "tidy"), *dirNameFlag).Output(); err != nil {
		panic(err)
	}
}

func setCommandDir(cmd *exec.Cmd, dir string) *exec.Cmd {
	cmd.Dir = dir
	return cmd
}

func writeMain() {
	content := fmt.Sprintf(`package main

import (
	"github.com/ahlixinjie/mongoose"
	"%s/internal/controller/grpc"
	"%s/internal/controller/http"
)

func main() {
	mongoose.Run(
		new(grpc.Impl),
		new(http.Impl),
	)
}
`, *moduleNamFlag, *moduleNamFlag)
	if err := os.WriteFile(fmt.Sprintf("%s/cmd/run/main.go", *dirNameFlag), []byte(content), 0644); err != nil {
		panic(err)
	}
}

func writeHttpController() {
	content := fmt.Sprintf(`package http

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	api "%s/api/http"
	"go.uber.org/dig"
	"google.golang.org/grpc"
)

type Impl struct {
}

func (i *Impl) Provide() (constructor interface{}, _ []dig.ProvideOption) {
	type conf struct {
		dig.Out
		Handler     *runtime.ServeMux
		GatewayFunc func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) (err error)
	}

	constructor = func() conf {
		return conf{
			Handler:     runtime.NewServeMux(),
			GatewayFunc: api.Register%sHandlerFromEndpoint,
		}
	}
	return
}
`, *moduleNamFlag, FirstUpper(serviceName()))
	if err := os.WriteFile(fmt.Sprintf("%s/internal/controller/http/init.go", *dirNameFlag), []byte(content), 0644); err != nil {
		panic(err)
	}
}

func writeGrpcController() {
	service := FirstUpper(serviceName())
	content := fmt.Sprintf(`package grpc

import (
	"context"
	api "%s/api/http"
	"go.uber.org/dig"
	"google.golang.org/grpc"
)

type Impl struct {
	api.Unimplemented%sServer
}

func (i *Impl) Invoke() (function interface{}, opts []dig.InvokeOption) {
	function = func(server *grpc.Server) {
		api.Register%sServer(server, new(Impl))
	}
	return
}

func (i *Impl) Echo(ctx context.Context, request *api.StringMessage) (response *api.StringMessage, err error) {
	response = &api.StringMessage{Value: request.GetValue()}
	return
}
`, *moduleNamFlag, service, service)
	if err := os.WriteFile(fmt.Sprintf("%s/internal/controller/grpc/init.go", *dirNameFlag), []byte(content), 0644); err != nil {
		panic(err)
	}
}

func writeServiceConfig() {
	content := `PORT:
  RPC: ":8080"
  HTTP: ":8081"
`
	if err := os.WriteFile(fmt.Sprintf("%s/config/dev.yaml", *dirNameFlag), []byte(content), 0644); err != nil {
		panic(err)
	}

	content = `PORT:
  RPC: "rpc"
  HTTP: "http"
`
	if err := os.WriteFile(fmt.Sprintf("%s/config/live.yaml", *dirNameFlag), []byte(content), 0644); err != nil {
		panic(err)
	}
}

func writeBufConfig() {
	apiConfigContent := fmt.Sprintf(`type: google.api.Service
config_version: 3

http:
  rules:
    - selector: %s.%s.Echo
      post: /v1/example/echo
      body: "*"`,
		protoPackageName(), serviceName())
	if err := os.WriteFile(fmt.Sprintf("%s/api/http/api.yaml", *dirNameFlag), []byte(apiConfigContent), 0644); err != nil {
		panic(err)
	}

	bufGenContent := `version: v1
plugins:
  - plugin: go
    out: .
    opt:
      - paths=source_relative
  - plugin: go-grpc
    out: .
    opt:
      - paths=source_relative
  - plugin: grpc-gateway
    out: .
    opt:
      - paths=source_relative
      - grpc_api_configuration=./api.yaml`
	if err := os.WriteFile(fmt.Sprintf("%s/api/http/buf.gen.yaml", *dirNameFlag), []byte(bufGenContent), 0644); err != nil {
		panic(err)
	}

	_, err := setCommandDir(exec.Command("buf", "generate"), fmt.Sprintf("%s/api/http", *dirNameFlag)).Output()
	if err != nil {
		panic(err)
	}

}

func writeProto() {
	content := fmt.Sprintf(`syntax = "proto3";
package %s;
option go_package = "%s/api/http";

message StringMessage {
  string value = 1;
}

service %s {
  rpc Echo(StringMessage) returns (StringMessage) {}
}`, protoPackageName(), *moduleNamFlag, serviceName())

	err := os.WriteFile(fmt.Sprintf("%s/api/http/service.proto", *dirNameFlag), []byte(content), 0644)
	if err != nil {
		panic(err)
	}
}

func protoPackageName() string {
	return strings.ReplaceAll(*moduleNamFlag, "/", ".") + ".api.http"
}

func serviceName() string {
	strs := strings.Split(*moduleNamFlag, "/")
	return strs[len(strs)-1]
}

// FirstUpper 字符串首字母大写
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	s = strings.ToLower(s)
	return strings.ToUpper(s[:1]) + s[1:]
}
