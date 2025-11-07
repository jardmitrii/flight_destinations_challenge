package main

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func fakeGetDestination(_, route string) (string, bool) {
	// Example: just return the route as destination if non-empty
	if route == "" {
		return "", false
	}
	return route, true
}

func TestWorker(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan job, 1)
	results := make(chan result, 1)
	defer close(results)

	// Prepare a job with routes including duplicates
	testJob := job{
		origin: "originTest",
		routes: []string{"dst1", "dst2", "dst1", "dst3", ""},
	}

	jobs <- testJob
	close(jobs)

	go worker(ctx, 1, jobs, results, fakeGetDestination)

	select {
	case res := <-results:
		expected := map[string]struct{}{
			"dst1": {},
			"dst2": {},
			"dst3": {},
		}
		if !reflect.DeepEqual(res.destinations, expected) {
			t.Errorf("unexpected destinations: got %v, want %v", res.destinations, expected)
		}
	case <-time.After(8 * time.Second):
		t.Error("timeout waiting for worker result")
	}
}

func TestWorker_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	jobs := make(chan job)
	results := make(chan result)
	defer close(jobs)

	// Cancel context immediately
	cancel()

	go func() {
		worker(ctx, 1, jobs, results, fakeGetDestination)
		close(results)
	}()

	select {
	case <-results:
		// Worker exited correctly on context cancel
	case <-time.After(1 * time.Second):
		t.Error("worker did not exit on context cancel promptly")
	}
}
