package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// начало решения

type Total struct {
	value atomic.Int32
}

func (total *Total) Increment() {
	total.value.Add(1)
}
func (total *Total) Value() int {
	return int(total.value.Load())
}

// конец решения

func main() {
	var wg sync.WaitGroup

	var total Total

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 10000; i++ {
				total.Increment()
			}
		}()
	}

	wg.Wait()
	fmt.Println("total", total.Value())
}
