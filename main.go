package main

import (
	"context"
	"fmt"
	"maps"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"
)

// flightRoutes - a predefined list of flight routes
// Each code is a concatenation of airport codes (e.g., LAXVIE = LAX to VIE)
// "VIELAX" could be treated as "LAXVIE", typically we have both destinations
var flightRoutes = []string{
	"LAXVIE",
	"laxfra",
	"BOSNIS",
	"LONVIE",
	"NYCLAX",
}

const (
	// chunkSize determines how many routes are processed in a single job
	chunkSize = 1
	// numberOfWorkers defines the concurrent worker threads processing routes
	numberOfWorkers = 1
)

func main() {
	// Handle graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Test demonstration destination counting for different origin airports
	fmt.Println(countDestinations(ctx, "LAX", flightRoutes))
	fmt.Println(countDestinations(ctx, "LON", flightRoutes))
	fmt.Println(countDestinations(ctx, "BOS", flightRoutes))

	os.Exit(0)
}

// countDestinations processes flight routes from a given origin
// Returns the number of unique destinations and a list of those destinations
func countDestinations(ctx context.Context, origin string, flightRoutes []string) (int, []string) {
	// Validate the airport code before processing
	err := validateAirportCode(origin)
	if err != nil {
		fmt.Println(err)
		return 0, nil
	}

	// Create a wait group to synchronize goroutines
	var wg sync.WaitGroup

	// Create channels for work and results
	jobs := make(chan job, numberOfWorkers)
	results := make(chan result, numberOfWorkers)

	// Spawn worker goroutines
	for i := 0; i < numberOfWorkers; i++ {
		wg.Go(func() { worker(ctx, i, jobs, results, getDestination) })
	}

	// Add jobs and close the jobs channel when done
	wg.Go(func() {
		defer close(jobs)
		addJobs(ctx, origin, flightRoutes, chunkSize, jobs)
	})

	// Cleanup
	go func() {
		// Wait for all work to be done
		wg.Wait()

		// Close results channel after all workers finish
		close(results)
	}()

	// Collect and deduplicate destinations
	res := map[string]struct{}{}
	for r := range results {
		fmt.Println("Results:", r)
		maps.Copy(res, r.destinations)
	}

	return len(res), slices.Collect(maps.Keys(res))
}
