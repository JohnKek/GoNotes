package main

import (
	"fmt"
	"reflect"
)

/*### Задача 1: Получение типа переменной

Описание: Напишите функцию, которая принимает переменную любого типа и возвращает ее тип в виде строки.

Подсказка: Используйте reflect.TypeOf() для получения типа переменной.
*/

// Функция для получения типа переменной
func GetType(v interface{}) string {
	return reflect.TypeOf(v).String()
}

func main() {
	// Примеры использования функции GetType
	fmt.Println(GetType(42))             // int
	fmt.Println(GetType("Hello"))        // string
	fmt.Println(GetType(3.14))           // float64
	fmt.Println(GetType(true))           // bool
	fmt.Println(GetType([]int{1, 2, 3})) // []int
}
