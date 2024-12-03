package main

import (
	"fmt"
	"reflect"
)

/*### Задача 3: Изменение значений полей структуры

Описание: Создайте структуру с несколькими полями. Напишите функцию, которая принимает указатель на эту структуру и изменяет значения ее полей на заданные.

Подсказка: Используйте reflect.Value для изменения значений полей, убедитесь, что вы работаете с указателем.
*/

type Config struct {
	Host string
	Port int
}

func UpdateConfig(c interface{}, host string, port int) {
	r := reflect.ValueOf(c).Elem()
	t := r.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == "Host" {
			r.Field(i).SetString(host)
		} else if field.Name == "Port" {
			r.Field(i).SetInt(int64(port))
		} else {
			fmt.Printf("Unknow field type: %s\n", field.Name)
		}
	}

}

func main() {
	cfg := Config{Host: "localhost", Port: 8080}
	fmt.Println(cfg.Host)
	fmt.Println(cfg.Port)
	UpdateConfig(&cfg, "localhost2", 8081)
	fmt.Println(cfg.Host)
	fmt.Println(cfg.Port)
}
