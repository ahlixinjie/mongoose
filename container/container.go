package container

import "go.uber.org/dig"

var container *dig.Container

func init() {
	container = dig.New()
}

func GetContainer() *dig.Container {
	return container
}
