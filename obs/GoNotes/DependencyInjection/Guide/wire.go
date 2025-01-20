// wire.go
//go:build wireinject
// +build wireinject

package main

import "github.com/google/wire"

// InitializeBar - функция для инициализации Bar
func InitializeBar() string {
	wire.Build(Set)
	return ""
}
