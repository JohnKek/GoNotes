package main

import (
	"fmt"
	"reflect"
)

/*
### Задача 2: Получение значений полей структуры

Описание: Создайте структуру с несколькими полями разных типов. Напишите функцию, которая принимает экземпляр этой структуры и выводит значения всех ее полей.

Подсказка: Используйте reflect.ValueOf() и NumField() для получения значений полей.
*/
type Person struct {
	Name string
	Age  int
}

func PrintFields(p interface{}) {
	v := reflect.ValueOf(p)
	fmt.Println(v)
	t := v.Type()
	fmt.Println(t)
	for i := 0; i < t.NumField(); i++ {
		fmt.Printf("%v:%v\n", t.Field(i).Name, v.Field(i))
	}
}

func main() {
	person := Person{Name: "Ivan", Age: 20}
	PrintFields(person)
}
