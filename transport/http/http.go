package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/ahlixinjie/mongoose/log"
	"github.com/ahlixinjie/mongoose/transport/common"
	"github.com/ahlixinjie/mongoose/utils/parse"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.uber.org/config"
	"go.uber.org/fx"
	"net/http"
	"os"
)

var (
	port int
	e    *echo.Echo
)

type params struct {
	fx.In
	Lc     fx.Lifecycle
	Config *config.YAML
}

func NewHTTPServer(p params) *echo.Echo {
	portStr := p.Config.Get(common.ConfKeyPort + "." + common.ConfKeyHTTP).String()
	port = parse.Port(portStr)
	if port == 0 && len(portStr) != 0 {
		//try to get port from env
		port = parse.Port(os.Getenv(portStr))
	}
	if port == 0 {
		log.GetLogger().Info("won't start http service")
		return nil
	}
	e = echo.New()
	e.Validator = &validatorWrapper{validator: validator.New()}

	p.Lc.Append(fx.StartHook(start))
	p.Lc.Append(fx.StopHook(stop))
	return e
}

func start() {
	go func() {
		log.GetLogger().Info("start http service")
		err := e.Start(fmt.Sprintf(":%d", port))
		if errors.Is(err, http.ErrServerClosed) {
			log.GetLogger().Info("http service has been shutdown")
			return
		}
		if err != nil {
			panic(err)
		}
	}()
}

func stop(ctx context.Context) error {
	log.GetLogger().Info("trying to stop http service")
	return e.Shutdown(ctx)
}
