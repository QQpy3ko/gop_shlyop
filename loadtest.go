package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	url := flag.String("url", "http://localhost:8080", "URL для тестирования")
	requests := flag.Int("n", 1000, "Количество запросов")
	concurrency := flag.Int("c", 50, "Количество одновременных горутин")
	flag.Parse()

	var wg sync.WaitGroup
	start := time.Now()

	// goroutines limit
	sem := make(chan struct{}, *concurrency)

	for i := 0; i < *requests; i++ {
		wg.Add(1)
		sem <- struct{}{} // take a slot

		go func() {
			defer wg.Done()
			resp, err := http.Get(*url)
			if err == nil {
				resp.Body.Close()
			}
			<-sem // release a slot
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	fmt.Printf("Выполнено %d запросов за %v\n", *requests, duration)
	fmt.Printf("Среднее время на запрос: %v\n", duration/time.Duration(*requests))
}
