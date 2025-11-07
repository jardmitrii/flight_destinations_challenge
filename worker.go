package main

import (
	"context"
	"fmt"
	"time"
)

// result contains destinations found during routs processing
type result struct {
	destinations map[string]struct{} // Map of destinations
}

// worker processes jobs from the jobs channel and sends results to the results channel
// Each worker handles a subset of routes concurrently
func worker(ctx context.Context, id int, jobs <-chan job, results chan<- result, work func(string, string) (string, bool)) {
	for {
		select {
		case j, ok := <-jobs:
			// Check if the jobs channel is closed and empty
			if !ok {
				return
			}

			// Debug log
			fmt.Printf("Worker %d processing job %q\n", id, j.routes)
			// Simulate work
			time.Sleep(time.Second * 5)

			// Create a map to store unique destinations
			destMap := map[string]struct{}{}

			for _, route := range j.routes {
				// Attempt to get destination for the current route
				dest, ok := work(j.origin, route)
				if !ok {
					// Skip routes that can't be processed
					continue
				}
				// Check if destination is already in map
				if _, ok := destMap[dest]; ok {
					continue
				}
				destMap[dest] = struct{}{}
			}
			// Send result with discovered destinations to results channel
			results <- result{destinations: destMap}
			// Debug log
			fmt.Printf("Worker %d finished job %q\n", id, j.routes)

		case <-ctx.Done():
			// Context canceled, stop working
			fmt.Printf("Worker %d shutting down\n", id)
			return
		}
	}
}
