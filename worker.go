package main

import (
	"context"
	"fmt"
	"sync"
)

// result contains destinations found during routs processing
type result struct {
	destinations map[string]struct{} // Map of destinations
}

// worker processes jobs from the jobs channel and sends results to the results channel
// Each worker handles a subset of routes concurrently
func startWorkers(ctx context.Context, workersCount int, jobs <-chan job, work func(job) result) <-chan result {
	// Create channel for results
	results := make(chan result, workersCount)
	// Create a wait group to synchronize goroutines
	var wg sync.WaitGroup

	// Spawn worker goroutines
	for i := 0; i < workersCount; i++ {
		wg.Go(func() {
			for {
				select {
				case j, ok := <-jobs:
					// Check if the jobs channel is closed and empty
					if !ok {
						return
					}

					// Debug log
					fmt.Printf("Worker %d processing job \n", i)
					r := work(j)
					// Debug log
					fmt.Printf("Worker %d finished job \n", i)

					select {
					case results <- r:
						fmt.Printf("Worker %d results ok\n", i)
					case <-ctx.Done():
						// Debug log
						fmt.Printf("Worker %d shutting down, results\n", i)
						return
						//default:
						//	results <- r
						//	fmt.Printf("Worker %d results ok\n", i)
					}

				case <-ctx.Done():
					// Context canceled, stop working
					fmt.Printf("Worker %d shutting down\n", i)
					return
				}
			}
		})
	}

	go func() {
		// Cleanup

		<-ctx.Done()
		// Wait for all work to be done
		wg.Wait()

		close(results)
	}()

	return results
}
