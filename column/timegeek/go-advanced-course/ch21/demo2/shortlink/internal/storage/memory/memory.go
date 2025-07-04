package memory

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/your_org/shortlink/internal/storage" // 导入定义的接口和类型
)

// 确保 *Store 实现了 storage.Store 接口 (编译时检查)
var _ storage.Store = (*Store)(nil)

// Store 是内存存储的具体实现类型
// 它实现了 storage.Store 接口
type Store struct {
	mu    sync.RWMutex
	links map[string]storage.Link // 使用 storage.Link
}

// NewStore 创建并返回一个新的内存存储实例
func New() *Store {
	return &Store{
		links: make(map[string]storage.Link),
	}
}

// Save 将 Link 对象保存到内存中
func (s *Store) Save(ctx context.Context, link storage.Link) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.links[link.ShortCode]; exists {
		return storage.ErrShortCodeExists // 使用 storage 包定义的错误
	}
	if link.CreatedAt.IsZero() {
		link.CreatedAt = time.Now().UTC()
	}
	s.links[link.ShortCode] = link
	log.Printf("DEBUG (memory.Store): Saved link. ShortCode: %s, LongURL: %s\n", link.ShortCode, link.LongURL)
	return nil
}

// FindByShortCode 根据短码从内存中查找 Link 对象
func (s *Store) FindByShortCode(ctx context.Context, shortCode string) (*storage.Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	link, exists := s.links[shortCode]
	if !exists {
		return nil, storage.ErrNotFound // 使用 storage 包定义的错误
	}
	linkCopy := link
	log.Printf("DEBUG (memory.Store): Found link. ShortCode: %s, LongURL: %s\n", link.ShortCode, link.LongURL)
	return &linkCopy, nil
}

// IncrementVisitCount 增加指定短码的访问计数
func (s *Store) IncrementVisitCount(ctx context.Context, shortCode string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	link, exists := s.links[shortCode]
	if !exists {
		return storage.ErrNotFound
	}
	link.VisitCount++
	s.links[shortCode] = link
	log.Printf("DEBUG (memory.Store): Incremented visit count. ShortCode: %s, NewCount: %d\n", link.ShortCode, link.VisitCount)
	return nil
}

// Close 对于内存存储，通常不需要特殊操作
func (s *Store) Close() error {
	log.Println("INFO (memory.Store): Closed (no-op)")
	return nil
}
