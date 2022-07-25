package mongoose

import (
	"fmt"
	"github.com/ahlixinjie/mongoose/container"
	"go.uber.org/dig"
	"testing"
)

type universe struct {
	p *plant
}
type plant struct {
	name string
	a    map[string]string
}

type air struct {
	name string
}

func TestDig(t *testing.T) {
	c := container.GetContainer()
	c.Provide(func() map[string]string {
		return map[string]string{"hi": "world"}
	}, dig.Name("conf"))

	var e plant

	type earthConf struct {
		dig.In
		A map[string]string `name:"conf"`
	}
	c.Provide(func(c earthConf) *plant {
		e.a = c.A
		return &e
	})
	fmt.Printf("%v", e.a)

	c.Invoke(func(p *plant) {
		u := universe{p: p}
		fmt.Println(u.p)
	})

}
