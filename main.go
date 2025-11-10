package main

import (
	"context"
	"fmt"
	"maps"
	"os/signal"
	"runtime"
	"slices"
	"syscall"
	"time"
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

	fmt.Println(runtime.NumGoroutine())
	// Add jobs and close the jobs channel when done
	jobs := addJobs(ctx, origin, flightRoutes, numberOfWorkers, chunkSize)
	results := startWorkers(ctx, numberOfWorkers, jobs, func(j job) result {
		// Debug log
		fmt.Printf("job %q\n", j)

		// Simulate work
		time.Sleep(time.Second * 1)

		// Create a map to store unique destinations
		destMap := map[string]struct{}{}

		for _, route := range j.routes {
			// Attempt to get destination for the current route
			dest, ok := getDestination(j.origin, route)
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

		return result{destinations: destMap}
	})

	// Collect and deduplicate destinations
	res := map[string]struct{}{}
	for r := range results {
		fmt.Println("Results:", r)
		maps.Copy(res, r.destinations)
	}

	return len(res), slices.Collect(maps.Keys(res))
}
