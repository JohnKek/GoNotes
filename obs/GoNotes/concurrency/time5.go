package main

import (
	"errors"
	"fmt"
	"time"
)

var ErrCanceled error = errors.New("canceled")

// начало решения

func withRateLimit(limit int, fn func()) (handle func() error, cancel func()) {
	canceled := make(chan struct{})
	interval := time.Second / time.Duration(limit)
	ticker := time.NewTicker(interval)
	handle = func() error {
		select {
		case <-ticker.C:
			go fn()
			return nil
		case <-canceled:
			return ErrCanceled
		}
	}
	cancel = func() {
		select {
		case <-canceled:
			return
		default:
			ticker.Stop()
			close(canceled)

		}
	}
	return handle, cancel
}

// конец решения

func main() {
	work := func() {
		fmt.Print(".")
	}

	handle, cancel := withRateLimit(5, work)
	defer cancel()

	start := time.Now()
	const n = 10
	for i := 0; i < n; i++ {
		handle()
	}
	fmt.Println()
	fmt.Printf("%d queries took %v\n", n, time.Since(start))
}
