package main

import "context"

// job is a unit of work containing an origin location and a slice of routes to process
type job struct {
	origin string   // Starting location for the routes
	routes []string // List of routes to be processed
}

// addJobs populates the jobs channel with chunks of routes to be processed
// It breaks down the full list of routes into smaller batches to manage workload
func addJobs(ctx context.Context, origin string, routes []string, chunkSize int, jobs chan<- job) {
	for i := 0; i < len(routes); i += chunkSize {
		select {
		case <-ctx.Done():
			// Stop adding jobs if context is canceled
			return
		default:
			// Calculate the end index for the current chunk, ensuring we don't exceed route list
			end := min(i+chunkSize, len(routes))
			// Send a job with a subset of routes to the jobs channel
			jobs <- job{origin: origin, routes: routes[i:end]}
		}
	}
}
