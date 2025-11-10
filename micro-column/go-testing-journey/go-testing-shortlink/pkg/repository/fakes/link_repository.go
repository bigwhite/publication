package fakes

import (
	"context"
	"sync"
	"time"

	"github.com/bigwhite/shortlink/pkg/domain"
	"github.com/bigwhite/shortlink/pkg/repository"
)

type FakeLinkRepository struct {
	mu    sync.RWMutex
	links map[string]*domain.Link
	idSeq int64
}

func NewFakeLinkRepository() *FakeLinkRepository {
	return &FakeLinkRepository{
		links: make(map[string]*domain.Link),
		idSeq: 1,
	}
}

func (f *FakeLinkRepository) FindByCode(ctx context.Context, code string) (*domain.Link, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	link, ok := f.links[code]
	if !ok {
		return nil, nil
	}
	// 返回一个副本，避免外部修改影响内部状态
	linkCopy := *link
	return &linkCopy, nil
}

func (f *FakeLinkRepository) Save(ctx context.Context, link *domain.Link) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 检查冲突
	if _, exists := f.links[link.ShortCode]; exists {
		return repository.ErrCodeConflict
	}

	// 模拟数据库行为
	link.ID = f.idSeq
	f.idSeq++
	link.CreatedAt = time.Now()

	f.links[link.ShortCode] = link
	return nil
}
