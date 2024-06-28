package main

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

func main() {
	var buf *bytes.Buffer
	//buf = new(bytes.Buffer)
	fmt.Println(buf)
	f(buf)
}
func f(out io.Writer) {
	// ... некоторые действия ... if out != nil {
	if out != nil {
		fmt.Println(1)
		out.Write([]byte("выполнено!\n"))
	}

}

func isEmpy(val interface{}) bool {
	return reflect.ValueOf(val).IsZero()
}

func setValue(structure interface{}, key string, value interface{}) {
	if reflect.TypeOf(structure).Elem().Kind() == reflect.Struct {
		elem := reflect.ValueOf(structure).Elem()
		field := elem.FieldByName(key)
		if field.CanSet() {
			field.Set(reflect.ValueOf(value))
		}
	}
}

func getValue(structure interface{}, key string) interface{} {
	var result interface{}
	if reflect.TypeOf(structure).Elem().Kind() == reflect.Struct {
		elem := reflect.ValueOf(structure).Elem()
		field := elem.FieldByName(key)
		if field.IsValid() {
			result = field.Interface()
		}
	}

	return result
}
