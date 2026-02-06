// fanout_fanin.go
//
// This example demonstrates an idiomatic Go concurrency pattern where a single
// producer fans out work to multiple goroutines, and their results are fanned in
// to a single results channel.
//
// Key ideas illustrated:
//
//   - Splitting independent work across goroutines (fan-out)
//   - Merging results into a single stream (fan-in)
//   - Using sync.WaitGroup to know when all workers are done
//   - Closing the results channel safely after workers finish
//
// This pattern is commonly used in real-world systems for:
//   - Parallel HTTP/API requests
//   - Database shard queries
//   - Scatter/gather workloads
//
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type job struct {
	id int
	n  int
}

type result struct {
	jobID int
	out   int
}

func worker(id int, jobs <-chan job, results chan<- result, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		// simulate variable work time
		time.Sleep(time.Duration(50+rand.Intn(150)) * time.Millisecond)
		results <- result{jobID: j.id, out: j.n * j.n}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	const (
		numWorkers = 4
		numJobs    = 10
	)

	jobs := make(chan job)
	results := make(chan result)

	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results, &wg)
	}

	// close results after all workers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// produce jobs
	go func() {
		for i := 1; i <= numJobs; i++ {
			jobs <- job{id: i, n: i}
		}
		close(jobs)
	}()

	for r := range results {
		fmt.Printf("job=%d square=%d\n", r.jobID, r.out)
	}
}
