package main

import (
	"fmt"
	"log"
	"reflect"
)

/*
### Задача 5: Создание экземпляра структуры динамически

Описание: Напишите функцию, которая принимает имя структуры и создает ее экземпляр с заданными значениями полей. Используйте reflect.New() для создания нового экземпляра.

Подсказка: Вам может понадобиться использовать reflect.TypeOf() для получения типа структуры.
*/
type User2 struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func CreateInstance(structName string, values ...interface{}) (interface{}, error) {
	var structType reflect.Type
	switch structName {
	case "User2":
		structType = reflect.TypeOf(User2{})
	default:
		return nil, fmt.Errorf("неизвестная структура: %s", structName)
	}

	instance := reflect.New(structType).Elem()

	for i := 0; i < instance.NumField(); i++ {
		if i < len(values) {
			instance.Field(i).Set(reflect.ValueOf(values[i]))
		}
	}

	return instance.Interface(), nil
}

func main() {
	instance, err := CreateInstance("User2", "user", "user")
	if err != nil {
		log.Printf("%s", err)
	}
	fmt.Println(instance)
}
