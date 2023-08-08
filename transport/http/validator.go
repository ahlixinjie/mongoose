package http

import (
	"github.com/ahlixinjie/mongoose/log"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type validatorWrapper struct {
	validator *validator.Validate
}

func (cv *validatorWrapper) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		log.GetLogger().Error("bad request", zap.Error(err))
		return status.Error(codes.InvalidArgument, "bad request")
	}
	return nil
}
