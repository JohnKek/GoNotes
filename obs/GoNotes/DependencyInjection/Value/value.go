// main.go
package main

import (
	"fmt"
)

// Foo - структура с полем X
type Foo struct {
	X int
}

func main() {
	// Инициализация Foo
	foo := injectFoo()
	fmt.Println(foo) // Вывод: {42}
}
