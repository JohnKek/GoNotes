package main

import (
	"fmt"
	"sync"
)

// начало решения

type Counter struct {
	mt *sync.Mutex
	kv map[string]int
}

func (c *Counter) Increment(str string) {
	c.mt.Lock()
	defer c.mt.Unlock()
	c.kv[str] += 1
	// ...
}

func (c *Counter) Value(str string) int {
	c.mt.Lock()
	defer c.mt.Unlock()
	return c.kv[str]
}

func (c *Counter) Range(fn func(key string, val int)) {
	c.mt.Lock()
	defer c.mt.Unlock()
	for key, value := range c.kv {
		fn(key, value)
	}

}

func NewCounter() *Counter {
	return &Counter{
		mt: &sync.Mutex{},
		kv: make(map[string]int),
	}

}

// конец решения

func main() {
	counter := NewCounter()

	var wg sync.WaitGroup
	wg.Add(3)

	increment := func(key string, val int) {
		defer wg.Done()
		for ; val > 0; val-- {
			counter.Increment(key)
		}
	}

	go increment("one", 100)
	go increment("two", 200)
	go increment("three", 300)

	wg.Wait()

	fmt.Println("two:", counter.Value("two"))

	fmt.Print("{ ")
	counter.Range(func(key string, val int) {
		fmt.Printf("%s:%d ", key, val)
	})
	fmt.Println("}")
}
