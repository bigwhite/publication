// pkg/repository/fakes/link_cache.go
package fakes

import (
	"context"
	"sync"
)

// FakeLinkCache 是 LinkCache 的一个内存实现，用于测试
type FakeLinkCache struct {
	mu     sync.RWMutex
	counts map[string]int64
}

func (f *FakeLinkCache) IncrementVisitCount(ctx context.Context, code string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.counts == nil {
		f.counts = make(map[string]int64)
	}
	f.counts[code]++
	return nil
}

func (f *FakeLinkCache) GetVisitCount(ctx context.Context, code string) (int64, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.counts[code], nil // 如果 key 不存在，返回 0，符合 Redis 行为
}
