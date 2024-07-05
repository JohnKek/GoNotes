package main

import (
	"fmt"
	"slices"
)

func main() {
	slice := []int{1, 2, 3, 4, 5}
	i := 2
	fmt.Println("slice:", slice)
	withAppend := append(slice[:i], slice[i+1:]...)
	fmt.Println("withAppend:", withAppend)
	fmt.Println("slice:", slice)
	letters := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	letters = slices.Delete(letters, 1, 1)
	fmt.Println("slice:", letters)
}
