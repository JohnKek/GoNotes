package main

import "fmt"

func main() {
	var Map map[string]string
	_, ok := Map["1"]
	Map["1"] = "1"
	fmt.Println(ok)
	fmt.Printf("%#v", len(Map))
}
