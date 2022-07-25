package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"go.uber.org/dig"
	"gopkg.in/yaml.v3"

	"github.com/ahlixinjie/mongoose/model"
	"github.com/ahlixinjie/mongoose/utils/env"
)

const (
	DigName = "dig_conf"
)

var F model.ProvideFunc = func() (constructor interface{}, opts []dig.ProvideOption) {
	e := env.GetEnv()

	file, err := os.Open(fmt.Sprintf("config/%s.yaml", e))
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	config := make(map[string]interface{})
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		panic(err)
	}
	return func() map[string]interface{} {
		return config
	}, []dig.ProvideOption{dig.Name(DigName)}
}
