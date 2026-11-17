package resilience

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestGateRespectsConcurrencyLimit(t *testing.T) {
	gate := NewGate(1)

	if err := gate.Acquire(context.Background()); err != nil {
		t.Fatalf("first acquire: %v", err)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		20*time.Millisecond,
	)
	defer cancel()

	err := gate.Acquire(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("second acquire error = %v, want deadline exceeded", err)
	}

	gate.Release()

	if err := gate.Acquire(context.Background()); err != nil {
		t.Fatalf("acquire after release: %v", err)
	}
	gate.Release()
}

func TestLimiter(t *testing.T) {
	limiter := NewLimiter(2, time.Second)

	if !limiter.Allow() {
		t.Fatal("first request should be allowed")
	}
	if !limiter.Allow() {
		t.Fatal("second request should be allowed")
	}
	if limiter.Allow() {
		t.Fatal("third request should be rejected")
	}
}
