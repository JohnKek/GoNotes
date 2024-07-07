package main

import "fmt"

// Bird базовый тип
type Bird struct{}

func (b *Bird) Fly() {
	fmt.Println("Птица летит")
}

// Penguin - подтип Bird, но не может летать
type Penguin struct {
	Bird
}

func main() {
	var bird = &Bird{}
	bird.Fly()

	var penguin = &Penguin{}
	penguin.Fly() // Нарушение LSP, т.к. пингвины не летают

}
