package main

import (
	"fmt"
	"github.com/google/wire"
)

type Foo int
type Bar int

func ProvideFoo() Foo { return 10 }

func ProvideBar() Bar { return 20 }

type FooBar struct {
	MyFoo Foo
	MyBar Bar
}

var Set = wire.NewSet(
	ProvideBar,
	ProvideFoo,
	wire.Struct(new(FooBar), "MyFoo", "MyBar"))

func main() {
	result := InitializeFooBar()
	fmt.Println(result)

}
