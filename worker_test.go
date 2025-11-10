package main

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func fakeGetDestination(j job) result {
	// Example: just return the routes as destinations if non-empty
	if j.routes == nil {
		return result{}
	}

	d := map[string]struct{}{}
	for _, route := range j.routes {
		if route == "" {
			continue
		}
		d[route] = struct{}{}
	}
	return result{destinations: d}
}

func Test_startWorkers(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobs := make(chan job, 1)

	// Prepare a job with routes including duplicates
	testJob := job{
		origin: "originTest",
		routes: []string{"dst1", "dst2", "dst1", "dst3", ""},
	}

	jobs <- testJob
	close(jobs)

	results := startWorkers(ctx, 1, jobs, fakeGetDestination)

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

func Test_startWorkers_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	jobs := make(chan job)
	defer close(jobs)

	// Cancel context immediately
	cancel()

	results := startWorkers(ctx, 1, jobs, fakeGetDestination)

	select {
	case <-results:
		// Worker exited correctly on context cancel
	case <-time.After(1 * time.Second):
		t.Error("worker did not exit on context cancel promptly")
	}
}
