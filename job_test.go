package main

import (
	"context"
	"reflect"
	"testing"
)

func TestAddJobs(t *testing.T) {
	tests := []struct {
		name      string
		origin    string
		routes    []string
		chunkSize int
		want      []job
	}{
		{
			name:      "normal case with chunk size 2",
			origin:    "origin1",
			routes:    []string{"r1", "r2", "r3", "r4", "r5"},
			chunkSize: 2,
			want: []job{
				{origin: "origin1", routes: []string{"r1", "r2"}},
				{origin: "origin1", routes: []string{"r3", "r4"}},
				{origin: "origin1", routes: []string{"r5"}},
			},
		},
		{
			name:      "empty routes",
			origin:    "origin2",
			routes:    []string{},
			chunkSize: 3,
			want:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			jobs := addJobs(ctx, tt.origin, tt.routes, len(tt.routes), tt.chunkSize)

			var got []job
			for j := range jobs {
				got = append(got, j)
			}
			cancel() // ensure cancel to avoid leaks

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addJobs() = %+v, want %+v", got, tt.want)
			}
			expectMaxJobs := 5
			if len(got) >= expectMaxJobs {
				t.Errorf("expected fewer jobs due to context cancel, but got %d jobs", len(got))
			}

		})
	}
}

func TestAddJobs_ContextCancelledBeforeAllJobsSent(t *testing.T) {
	origin := "originTest"
	routes := []string{"r1", "r2", "r3", "r4", "r5", "r6"}
	chunkSize := 2

	ctx, cancel := context.WithCancel(context.Background())
	jobs := addJobs(ctx, origin, routes, len(routes)/chunkSize+1, chunkSize)

	// Cancel context to simulate early termination
	cancel()

	var got []job
	for j := range jobs {
		got = append(got, j)
	}

	expectMaxJobs := (len(routes) + chunkSize - 1) / chunkSize
	if len(got) >= expectMaxJobs {
		t.Errorf("expected fewer jobs due to context cancel, but got %d jobs", len(got))
	}
}
