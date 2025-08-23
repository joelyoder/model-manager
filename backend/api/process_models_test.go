package api

import (
	"sync"
	"testing"
	"time"

	"bou.ke/monkey"
)

func TestProcessModelsConcurrency(t *testing.T) {
	// prepare more than 4 items
	items := make([]CivitModel, 5)

	started := make(chan struct{}, len(items))
	release := make(chan struct{})

	var mu sync.Mutex
	current := 0
	maxConc := 0

	patch := monkey.Patch(processModel, func(item CivitModel, apiKey string) {
		mu.Lock()
		current++
		if current > maxConc {
			maxConc = current
		}
		mu.Unlock()

		started <- struct{}{}
		<-release

		mu.Lock()
		current--
		mu.Unlock()
	})
	defer patch.Unpatch()

	done := make(chan struct{})
	go func() {
		processModels(items, "")
		close(done)
	}()

	// Wait for first four to start
	for i := 0; i < 4; i++ {
		<-started
	}

	// Ensure the fifth has not started yet
	select {
	case <-started:
		t.Fatalf("more than four models processed concurrently")
	case <-time.After(50 * time.Millisecond):
	}

	// Release the first four
	for i := 0; i < 4; i++ {
		release <- struct{}{}
	}

	// Now the fifth should start
	<-started
	release <- struct{}{}

	<-done

	mu.Lock()
	got := maxConc
	mu.Unlock()

	if got > 4 {
		t.Fatalf("expected at most 4 concurrent processes, got %d", got)
	}
	if got < 4 {
		t.Fatalf("expected concurrency limit to reach 4, got %d", got)
	}
}
