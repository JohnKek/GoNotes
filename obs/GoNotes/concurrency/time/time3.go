package main

import (
	"fmt"
	"math/rand"
	"time"
)

// начало решения

func delay(dur time.Duration, fn func()) func() {
	timer := time.NewTimer(dur)
	cancel := make(chan struct{})

	go func() {
		select {
		case <-timer.C:
			fn()
		case <-cancel:
			fmt.Println(1)
			return
		}
	}()

	cancelFn := func() {
		if timer.Stop() {
			close(cancel)
		}
	}

	return cancelFn
}

// конец решения

func main() {
	rand.Seed(time.Now().Unix())

	work := func() {
		fmt.Println("work done")
	}

	cancel := delay(100*time.Millisecond, work)

	time.Sleep(10 * time.Millisecond)
	if rand.Float32() < 0.5 {
		cancel()
		fmt.Println("delayed function canceled")
	}
	time.Sleep(100 * time.Millisecond)
}
