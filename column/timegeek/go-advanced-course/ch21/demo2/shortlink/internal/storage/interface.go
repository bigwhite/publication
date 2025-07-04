package storage

import (
	"context"
	"errors"
	"time"
)

// Link 是存储在数据库中的短链接信息。
// 它的字段都应可导出，以便 Service 层和可能的其他包使用。
type Link struct {
	ShortCode  string    // 短码，唯一标识
	LongURL    string    // 原始长链接
	VisitCount int64     // 访问次数
	CreatedAt  time.Time // 创建时间
}

// Store 定义了数据存储层需要提供的核心能力。
// 所有实现都应该是并发安全的。
type Store interface {
	Save(ctx context.Context, link Link) error
	FindByShortCode(ctx context.Context, shortCode string) (*Link, error)
	IncrementVisitCount(ctx context.Context, shortCode string) error
	Close() error
}

// 包级别导出的哨兵错误，供调用者使用 errors.Is 判断。
var ErrNotFound = errors.New("storage: link not found")
var ErrShortCodeExists = errors.New("storage: short code already exists")
