package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
)

func main() {
	t := reflect.TypeOf(3)  // reflect.Type
	fmt.Println(t.String()) // "int"
	fmt.Println(t)
	var w io.Writer = os.Stdout
	fmt.Println(reflect.TypeOf(w))
}
