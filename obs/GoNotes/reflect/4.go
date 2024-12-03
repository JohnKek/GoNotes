package main

import (
	"fmt"
	"reflect"
)

/*
### Задача 4: Проверка наличия тега в полях структуры

Описание: Создайте структуру с полями, у которых есть теги. Напишите функцию, которая принимает экземпляр этой структуры и выводит значения полей вместе с их тегами.

Подсказка: Используйте reflect.StructTag для получения тегов полей.
*/
type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func PrintFieldsTags(p interface{}) {
	v := reflect.ValueOf(p)
	fmt.Println(v)
	t := v.Type()
	fmt.Println(t)
	for i := 0; i < t.NumField(); i++ {
		fmt.Printf("%v:%v:%v\n", t.Field(i).Name, v.Field(i), t.Field(i).Tag)
	}
}

func main() {
	u := User{
		Username: "username",
		Email:    "email",
	}
	PrintFieldsTags(u)
}
