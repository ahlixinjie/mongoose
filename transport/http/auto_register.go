package http

import (
	"context"
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

func AutoRegister(e *echo.Echo, server interface{}, prefix string, code2Http map[codes.Code]int) {
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
	if !methodType.Out(1).Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
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
	methods map[string]*methodWrap, code2Http map[codes.Code]int) {
	codec := defaultCodec{code2Http: code2Http}
	for methodName, method := range methods {
		e.POST(prefix+methodName, func(c echo.Context) error {
			req := reflect.New(method.ReqType).Interface()
			err := codec.Decode(c.Request(), req)
			if err != nil {
				return err
			}

			ctx, err := GenCtxFromEcho(c)
			if err != nil {
				return err
			}

			res := method.method.Func.Call([]reflect.Value{receiverValue, reflect.ValueOf(ctx), reflect.ValueOf(req)})
			return codec.EncodeResponse(c.Request(), c.Response(), res[0].Interface(), res[1].Interface().(error))
		})
	}
}

func GenCtxFromEcho(c echo.Context) (context.Context, error) {
	var md = make(metadata.MD)
	if err := (&echo.DefaultBinder{}).BindHeaders(c, &md); err != nil {
		return nil, err
	}

	if len(md.Get(common.HeaderRequestID)) == 0 {
		md.Set(common.HeaderRequestID, strings.Join(strings.Split(uuid.New().String(), "-"), ""))
	}
	ctx := metadata.NewIncomingContext(c.Request().Context(), md)
	c.SetRequest(c.Request().WithContext(ctx))
	return ctx, nil
}
