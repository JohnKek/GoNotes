// wire.go
//go:build wireinject
// +build wireinject

package main

import "github.com/google/wire"

// injectFoo - функция для внедрения зависимости Foo
func injectFoo() Foo {
	wire.Build(wire.Value(Foo{X: 42}))
	return Foo{}
}
