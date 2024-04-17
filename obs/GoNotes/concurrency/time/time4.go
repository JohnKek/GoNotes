package main

import (
	"fmt"
	"time"
)

// начало решения

func schedule(dur time.Duration, fn func()) func() {
	cancel := make(chan struct{})
	ticker := time.NewTicker(dur)
	go func() {
		for {
			select {
			case <-ticker.C:
				fn()
			case <-cancel:
				return
			}

		}

	}()
	return func() {
		select {
		case <-cancel:
			fmt.Println(1)
			return
		default:
			ticker.Stop()
			close(cancel)
			fmt.Println(0)
		}
	}
}

// конец решения

func main() {
	work := func() {
		at := time.Now()
		fmt.Printf("%s: work done\n", at.Format("15:04:05.000"))
	}

	cancel := schedule(50*time.Millisecond, work)
	time.Sleep(260 * time.Millisecond)
	cancel()
	cancel()
}
