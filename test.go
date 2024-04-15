package GoNotes

import "sync"

func merge(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	size := len(channels)
	wg.Add(size)
	out := make(chan int)
	for _, chanel := range channels {
		go func(ch <-chan int) {
			defer wg.Done()
			for val := range chanel {
				out <- val
			}
		}(chanel)
	}
	go func() {
		defer close(out)
		wg.Wait()
	}()
	return out
	// объедините все исходные каналы в один выходной
	// последовательное объединение НЕ подходит
}
