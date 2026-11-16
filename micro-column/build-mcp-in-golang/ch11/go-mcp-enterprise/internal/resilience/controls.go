package resilience

import (
	"context"
	"sync"
	"time"
)

// Gate 限制同时执行的操作数量。
type Gate struct {
	sem chan struct{}
}

func NewGate(maxConcurrent int) *Gate {
	if maxConcurrent <= 0 {
		panic("maxConcurrent must be greater than zero")
	}

	return &Gate{
		sem: make(chan struct{}, maxConcurrent),
	}
}

func (g *Gate) Acquire(ctx context.Context) error {
	select {
	case g.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (g *Gate) Release() {
	<-g.sem
}

// Limiter 是一个简单的固定窗口限流器。
// 它适合单进程演示，不提供跨实例的全局限流能力。
type Limiter struct {
	mu        sync.Mutex
	limit     int
	window    time.Duration
	startedAt time.Time
	count     int
}

func NewLimiter(limit int, window time.Duration) *Limiter {
	if limit <= 0 {
		panic("limit must be greater than zero")
	}
	if window <= 0 {
		panic("window must be greater than zero")
	}

	return &Limiter{
		limit:  limit,
		window: window,
	}
}

func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	if l.startedAt.IsZero() || now.Sub(l.startedAt) >= l.window {
		l.startedAt = now
		l.count = 0
	}

	if l.count >= l.limit {
		return false
	}

	l.count++
	return true
}
