//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import "github.com/google/wire"

// InitializeEvent creates an Event. It will error if the Event is staffed with
// a grumpy greeter.
func InitializeFooBar() FooBar {
	wire.Build(Set)
	return FooBar{}
}
