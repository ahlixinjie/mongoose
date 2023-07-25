package http

import (
	"context"
	stringconverter "github.com/ahlixinjie/go-utils/string/converter"
	"github.com/ahlixinjie/mongoose/transport/common"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"reflect"
	"strings"
)

type methodWrap struct {
	ReqType  reflect.Type
	RespType reflect.Type
	method   reflect.Method
}

func AutoRegister(e *echo.Echo, server interface{}, prefix string, code2Http func(codes.Code) int) {
	receiverValue := reflect.ValueOf(server)
	receiverType := reflect.TypeOf(server)
	methods := make(map[string]*methodWrap)

	for i := 0; i < receiverType.NumMethod(); i++ {
		method := receiverType.Method(i)
		if m := validateMethod(method); m != nil {
			methods[method.Name] = m
		}
	}
	registerMethod(receiverValue, e, prefix, methods, code2Http)
}

func validateMethod(method reflect.Method) *methodWrap {
	if !method.IsExported() {
		return nil
	}

	methodType := method.Type
	numIn, numOut := methodType.NumIn(), methodType.NumOut()
	if numIn != 3 || numOut != 2 {
		return nil
	}

	//first arg needs to be context
	if !methodType.In(1).Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
		return nil
	}

	//last return arg needs to be error
	if !methodType.Out(1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return nil
	}

	reqType := methodType.In(2)
	if reqType.Kind() != reflect.Ptr {
		return nil
	}

	respType := methodType.In(0)
	if respType.Kind() != reflect.Ptr {
		return nil
	}

	return &methodWrap{
		ReqType:  methodType.In(2).Elem(),
		RespType: methodType.Out(0).Elem(),
		method:   method,
	}
}

func registerMethod(receiverValue reflect.Value, e *echo.Echo, prefix string,
	methods map[string]*methodWrap, code2Http func(codes.Code) int) {
	codec := defaultCodec{code2Http: code2Http}
	for k, v := range methods {
		methodName, method := k, v
		e.POST(prefix+stringconverter.CamelCaseToUnderscore(methodName), func(c echo.Context) error {
			reply, respErr := func() (interface{}, error) {
				req := reflect.New(method.ReqType).Interface()
				err := codec.Decode(c.Request(), req)
				if err != nil {
					return nil, err
				}
				res := method.method.Func.Call([]reflect.Value{receiverValue, reflect.ValueOf(GenCtxFromEcho(c)), reflect.ValueOf(req)})
				if res[1].Interface() != nil {
					return nil, res[1].Interface().(error)
				}
				return res[0].Interface(), nil
			}()
			return codec.EncodeResponse(c.Request(), c.Response(), reply, respErr)
		})
	}
}

func GenCtxFromEcho(c echo.Context) context.Context {
	var md = metadata.New(make(map[string]string))
	for k, v := range c.Request().Header.Clone() {
		md.Append(k, v...)
	}

	if len(md.Get(common.HeaderRequestID)) == 0 {
		md.Set(common.HeaderRequestID, strings.Join(strings.Split(uuid.New().String(), "-"), ""))
	}
	ctx := metadata.NewIncomingContext(c.Request().Context(), md)
	c.SetRequest(c.Request().WithContext(ctx))
	return ctx
}
