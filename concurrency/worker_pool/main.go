// worker_pool.go
//
// This example demonstrates a worker pool that processes jobs concurrently.
//
// Key ideas illustrated:
//
//   - Fixed number of workers consuming a jobs channel
//   - Collecting results via a results channel
//   - Coordinating completion with sync.WaitGroup
//
package main

import (
	"fmt"
	"sync"
)

type Job struct {
	ID int
}

type Result struct {
	ID   int
	Text string
}

func worker(jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		results <- Result{ID: j.ID, Text: fmt.Sprintf("processed job %d", j.ID)}
	}
}

func main() {
	const workers = 3
	const total = 8

	jobs := make(chan Job)
	results := make(chan Result)

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker(jobs, results, &wg)
	}

	go func() {
		for i := 1; i <= total; i++ {
			jobs <- Job{ID: i}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		fmt.Println(r.Text)
	}
}
