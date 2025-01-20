// main.go
package main

import (
	"fmt"
	"github.com/google/wire"
)

// Fooer - интерфейс с методом Foo
type Fooer interface {
	Foo() string
}

// MyFooer - структура, реализующая интерфейс Fooer
type MyFooer string

// Foo - метод, реализующий интерфейс Fooer
func (b *MyFooer) Foo() string {
	return string(*b)
}

// provideMyFooer - функция для создания MyFooer
func provideMyFooer() *MyFooer {
	b := new(MyFooer)
	*b = "Hello, World!"
	return b
}

// Bar - структура, которая будет использовать Fooer
type Bar string

// provideBar - функция, которая принимает Fooer и возвращает строку
func provideBar(f Fooer) string {
	// f будет *MyFooer.
	return f.Foo()
}

// Set - набор зависимостей для wire
var Set = wire.NewSet(
	provideMyFooer,
	wire.Bind(new(Fooer), new(*MyFooer)),
	provideBar,
)

func main() {
	// Инициализация Bar
	result := InitializeBar()
	fmt.Println(result) // Вывод: Hello, World!
}
