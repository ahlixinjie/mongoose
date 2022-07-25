package model

import "go.uber.org/dig"

type Provider interface {
	Provide() (constructor interface{}, opts []dig.ProvideOption)
}

type ProvideFunc func() (constructor interface{}, opts []dig.ProvideOption)

func (p ProvideFunc) Provide() (constructor interface{}, opts []dig.ProvideOption) {
	return p()
}

type Invoker interface {
	Invoke() (function interface{}, opts []dig.InvokeOption)
}

type StarterConf struct {
	S        Starter
	Priority uint
}

type StopperConf struct {
	S        Stopper
	Priority uint
}

type Starter interface {
	Start() error
}

type Stopper interface {
	Stop() error
}

type StarterPriority interface {
	Starter
	Priority() uint
}

type StopperPriority interface {
	Stopper
	Priority() uint
}
