package main

import (
	"context"
	"fmt"
	"time"
)

// job is a unit of work containing an origin location and a slice of routes to process
type job struct {
	origin string   // Starting location for the routes
	routes []string // List of routes to be processed
}

// result contains destinations found during routs processing
type result struct {
	destinations map[string]struct{} // Map of destinations
}

// addJobs populates the jobs channel with chunks of routes to be processed
// It breaks down the full list of routes into smaller batches to manage workload
func addJobs(ctx context.Context, origin string, routes []string, chunkSize int, jobs chan<- job) {
	for i := 0; i < len(flightRoutes); i += chunkSize {
		select {
		case <-ctx.Done():
			// Stop adding jobs if context is canceled
			return
		default:
			// Calculate the end index for the current chunk, ensuring we don't exceed route list
			end := min(i+chunkSize, len(routes))
			// Send a job with a subset of routes to the jobs channel
			jobs <- job{origin: origin, routes: flightRoutes[i:end]}
		}
	}
}

// worker processes jobs from the jobs channel and sends results to the results channel
// Each worker handles a subset of routes concurrently
func worker(ctx context.Context, id int, jobs <-chan job, results chan<- result) {
	for {
		select {
		case job, ok := <-jobs:
			// Check if the jobs channel is closed and empty
			if !ok {
				return
			}

			// Debug log
			fmt.Printf("Worker %d processing job %q\n", id, job.routes)
			// Simulate work
			time.Sleep(time.Second * 5)

			// Create a map to store unique destinations
			destMap := map[string]struct{}{}

			for _, route := range job.routes {
				// Attempt to get destination for the current route
				dest, ok := getDestination(job.origin, route)
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
			fmt.Printf("Worker %d finished job %q\n", id, job.routes)

		case <-ctx.Done():
			// Context canceled, stop working
			fmt.Printf("Worker %d shutting down\n", id)
			return
		}
	}
}
