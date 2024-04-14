package main

import "fmt"

func multyply(x int) (y int) {
	defer func() { y = y * 2 }()
	return x * x
}
func main() {
	X := 5
	fmt.Println(multyply(X))
}
