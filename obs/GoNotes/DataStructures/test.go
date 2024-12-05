package main

import (
	"fmt"
	"sync/atomic"
)

func main() {
	a := atomic.Uint32{}
	a.Store(2)
	fmt.Println(a.Load())
	a.Add(^uint32(0))
	fmt.Println(a.Load())
}
