package concurrency

import (
	"sync"
	"testing"
)

var counter int // Shared variable that will cause a data race

// incrementCounter increments the global counter. This is not safe for concurrent use.
func incrementCounter(wg *sync.WaitGroup) {
	defer wg.Done()
	// This line is the source of the data race when called concurrently:
	counter++
}

// TestDataRace demonstrates a data race condition.
// Run with `go test -race -run TestDataRace` to detect it.
func TestDataRace(t *testing.T) {
	var wg sync.WaitGroup
	numGoroutines := 100
	counter = 0 // Reset counter for each test run

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go incrementCounter(&wg)
	}
	wg.Wait() // Wait for all goroutines to complete

	// The final value of `counter` is non-deterministic due to the race.
	// We don't assert its value because it's unpredictable.
	// The primary purpose of this test (when run with -race) is for the
	// race detector to report the concurrent unprotected access to `counter`.
	t.Logf("Counter value after concurrent increments (non-deterministic due to race): %d", counter)
	// A simple assertion just to make the test do something checkable,
	// though the real check is the race detector's output.
	if numGoroutines > 0 && counter < 1 { // This condition is arbitrary for demo
		// t.Errorf("Counter should be greater than 0 if goroutines ran, but this is not a reliable check without -race.")
	}
}

// --- Corrected version with a Mutex to prevent data race ---
var safeCounter int
var mu sync.Mutex // Mutex to protect safeCounter

// incrementSafeCounter increments the global safeCounter using a mutex.
func incrementSafeCounter(wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock() // Acquire lock
	safeCounter++
	mu.Unlock() // Release lock
}

// TestNoDataRace demonstrates the correct way to handle shared state concurrently.
// Run with `go test -race -run TestNoDataRace` to verify no race is detected.
func TestNoDataRace(t *testing.T) {
	var wg sync.WaitGroup
	numGoroutines := 100
	safeCounter = 0 // Reset safeCounter

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go incrementSafeCounter(&wg)
	}
	wg.Wait()

	// With the mutex, the final value of safeCounter should be deterministic.
	if safeCounter != numGoroutines {
		t.Errorf("Expected safeCounter to be %d, got %d", numGoroutines, safeCounter)
	} else {
		t.Logf("SafeCounter correctly incremented to: %d", safeCounter)
	}
}
