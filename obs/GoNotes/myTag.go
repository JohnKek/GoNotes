package main

import (
	"fmt"
	"reflect"
)

type User struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}

func ToMap(v interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		tag := field.Tag.Get("db")
		if tag != "" {
			result[tag] = val.Field(i).Interface()
		}
	}

	return result
}

func main() {
	user := User{ID: 1, Name: "Alice", Age: 30}
	userMap := ToMap(&user)

	fmt.Println(userMap) // Вывод: map[age:30 id:1 name:Alice]
}
