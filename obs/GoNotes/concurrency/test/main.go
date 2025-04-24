package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// if ok print ""ok {url}" else "error {url}"
func main() {
	urls := []string{"many urls..."} // 1 000 000
	workerCount := 10
	var wg sync.WaitGroup
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	urlsChan := make(chan string, 10)
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(ctx, urlsChan, &wg)
	}

	for _, urlElem := range urls {
		urlsChan <- urlElem
	}

	close(urlsChan)

	wg.Wait()
}

func fetch(url string) (int, error) {
	// fetch to url
	code, err := http.Get(url)
	return code, err
}

func worker(ctx context.Context, urls chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case url, ok := <-urls:
			if !ok {
				return
			}
			code, err := fetch(url)
			if err != nil && code == 200 {
				fmt.Sprintf("ok %s", url)
				break
			}
			fmt.Sprintf("err %s", err)
		}
	}
}
