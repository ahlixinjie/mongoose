package mongoose

import (
	"bytes"
	"fmt"
	"github.com/ahlixinjie/mongoose/config"
	"github.com/ahlixinjie/mongoose/container"
	"github.com/ahlixinjie/mongoose/model"
	"github.com/ahlixinjie/mongoose/transport/grpc"
	"go.uber.org/dig"
	"os"
	"os/signal"
	"sort"
	"syscall"
)

type app struct {
	providers []model.Provider
	invokers  []model.Invoker
	starters  []model.StarterConf
	stoppers  []model.StopperConf

	exitChan chan struct{}
}

func Run(elements ...interface{}) {
	elements = append(elements, config.F, &grpc.Service{})
	a := newApp(elements...)
	a.run()
	a.waitStop()
}

func (a *app) run() {
	defer func() {
		b := &bytes.Buffer{}
		if err := dig.Visualize(container.GetContainer(), b); err != nil {
			panic(err)
		}
		os.WriteFile("./app.dot", b.Bytes(), 0644)
	}()
	for _, v := range a.providers {
		c, opts := v.Provide()
		if err := container.GetContainer().Provide(c, opts...); err != nil {
			panic(err)
		}
	}
	for _, v := range a.invokers {
		f, opts := v.Invoke()
		if err := container.GetContainer().Invoke(f, opts...); err != nil {
			panic(err)
		}
	}

	for _, v := range a.starters {
		if err := v.S.Start(); err != nil {
			panic(err)
		}
	}
}

func (a *app) waitStop() {
	a.exitChan = make(chan struct{})
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		select {
		case v := <-c:
			fmt.Printf("received %v\n", v)
			a.exitChan <- struct{}{}
		}
	}()

	<-a.exitChan

	for _, v := range a.stoppers {
		if err := v.S.Stop(); err != nil {
			panic(err)
		}
	}
	fmt.Println("exit")
}

func newApp(elements ...interface{}) *app {
	a := app{}
	for _, e := range elements {
		if p, ok := e.(model.Provider); ok {
			a.providers = append(a.providers, p)
		}
		if i, ok := e.(model.Invoker); ok {
			a.invokers = append(a.invokers, i)
		}

		if s, ok := e.(model.StarterPriority); ok {
			a.starters = append(a.starters, model.StarterConf{
				S:        s,
				Priority: s.Priority(),
			})
		} else if ss, ok := e.(model.Starter); ok {
			a.starters = append(a.starters, model.StarterConf{
				S:        ss,
				Priority: 0,
			})
		}

		if s, ok := e.(model.StopperPriority); ok {
			a.stoppers = append(a.stoppers, model.StopperConf{
				S:        s,
				Priority: s.Priority(),
			})
		} else if ss, ok := e.(model.Stopper); ok {
			a.stoppers = append(a.stoppers, model.StopperConf{
				S:        ss,
				Priority: 0,
			})
		}
	}
	sort.SliceStable(a.starters, func(i, j int) bool {
		return a.starters[i].Priority > a.starters[j].Priority
	})
	sort.SliceStable(a.stoppers, func(i, j int) bool {
		return a.stoppers[i].Priority > a.stoppers[j].Priority
	})
	return &a
}
