package memory

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/your_org/shortlink/internal/store" // 导入定义的接口和错误
)

// MemoryStore 是一个基于内存的Store接口实现，主要用于测试和演示
type MemoryStore struct {
	mu     sync.RWMutex
	links  map[string]*store.LinkEntry // shortCode -> LinkEntry
	logger *slog.Logger
}

// NewStore 创建一个新的MemoryStore实例
func NewStore(logger *slog.Logger) store.Store {
	return &MemoryStore{
		links:  make(map[string]*store.LinkEntry),
		logger: logger.With("component", "memory_store"),
	}
}

// Save 实现Store接口的Save方法
func (s *MemoryStore) Save(ctx context.Context, entry *store.LinkEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查短码是否已存在 (在真实DB中，这通常由唯一约束保证)
	if _, exists := s.links[entry.ShortCode]; exists {
		s.logger.WarnContext(ctx, "Attempted to save duplicate short code",
			slog.String("short_code", entry.ShortCode))
		return store.ErrConflict
	}

	// 简单复制一下，避免外部修改影响内部存储 (虽然对于内存存储可能不是大问题)
	entryToSave := *entry
	if entryToSave.CreatedAt.IsZero() {
		entryToSave.CreatedAt = time.Now()
	}
	s.links[entry.ShortCode] = &entryToSave

	s.logger.InfoContext(ctx, "Link entry saved to memory store",
		slog.String("short_code", entry.ShortCode),
		slog.String("long_url", entry.LongURL),
	)
	return nil
}

// FindByShortCode 实现Store接口的FindByShortCode方法
func (s *MemoryStore) FindByShortCode(ctx context.Context, shortCode string) (*store.LinkEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.links[shortCode]
	if !exists {
		s.logger.DebugContext(ctx, "Short code not found in memory store", slog.String("short_code", shortCode))
		return nil, store.ErrNotFound
	}

	// 检查是否过期 (如果实现了过期逻辑)
	if !entry.ExpireAt.IsZero() && time.Now().After(entry.ExpireAt) {
		s.logger.InfoContext(ctx, "Short code found but has expired", slog.String("short_code", shortCode))
		// (可选) 可以在这里从map中删除过期的条目
		// delete(s.links, shortCode) // 需要 s.mu.Lock()
		return nil, store.ErrNotFound // 对外表现为未找到
	}

	s.logger.DebugContext(ctx, "Link entry found in memory store", slog.String("short_code", shortCode))
	// 返回一个副本，避免外部修改
	retEntry := *entry
	return &retEntry, nil
}

// Close 实现Store接口的Close方法 (对于内存存储，通常无事可做)
func (s *MemoryStore) Close() error {
	s.logger.Info("Memory store closing (no-op).")
	return nil
}

// (其他Store接口方法的简单实现或占位)
// func (s *MemoryStore) FindByOriginalURLAndUserID(...) (*store.LinkEntry, error) { ... return nil, store.ErrNotFound }
// func (s *MemoryStore) IncrementVisitCount(...) error { ... return nil }
// func (s *MemoryStore) GetVisitCount(...) (int64, error) { ... return 0, nil }
