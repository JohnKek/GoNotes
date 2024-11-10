```go
package main

import (
	"fmt"
)

// Clonable - интерфейс для объектов, поддерживающих клонирование
type Clonable interface {
	Clone() Clonable
}

// Address - структура, представляющая адрес
type Address struct {
	City  string
	State string
}

// Clone - метод для клонирования объекта Address
func (a *Address) Clone() Clonable {
	return &Address{
		City:  a.City,
		State: a.State,
	}
}

// Person - структура, представляющая человека
type Person struct {
	Name    string
	Age     int
	Address *Address // Поле, которое также нужно клонировать
}

// Clone - метод для клонирования объекта Person
func (p *Person) Clone() Clonable {
	return &Person{
		Name:    p.Name,
		Age:     p.Age,
		Address: p.Address.Clone().(*Address), // Клонируем адрес
	}
}

func main() {
	// Создаем оригинальный объект Person с адресом
	original := &Person{
		Name: "Alice",
		Age:  30,
		Address: &Address{
			City:  "New York",
			State: "NY",
		},
	}

	// Клонируем объект Person
	clone := original.Clone().(*Person)

	// Изменяем клонированный объект и его адрес
	clone.Name = "Bob"
	clone.Age = 25
	clone.Address.City = "Los Angeles"
	clone.Address.State = "CA"

	// Выводим оригинальный и клонированный объекты
	fmt.Printf("Original: %+v\n", original)
	fmt.Printf("Clone: %+v\n", clone)
}

```