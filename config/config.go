package config

import (
	"fmt"
	"github.com/ahlixinjie/mongoose/utils/env"
	"go.uber.org/config"
)

const (
	InjectName = "conf"
)

func NewProvider() (*config.YAML, error) {
	yaml, err := config.NewYAML(config.File(fmt.Sprintf("config/%s.yaml", env.GetEnv())))
	if err != nil {
		return nil, err
	}
	return yaml, nil
}
