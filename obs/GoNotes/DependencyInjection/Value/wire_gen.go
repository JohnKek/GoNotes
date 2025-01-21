// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

// Injectors from wire.go:

// injectFoo - функция для внедрения зависимости Foo
func injectFoo() Foo {
	foo := _wireFooValue
	return foo
}

var (
	_wireFooValue = Foo{X: 42}
)
