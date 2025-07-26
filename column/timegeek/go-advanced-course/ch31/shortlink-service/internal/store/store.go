// internal/store/store.go
package store

import (
	"context"
	"errors"
	"time"
)

// ErrNotFound 表示请求的资源未找到
var ErrNotFound = errors.New("store: requested item not found")

// ErrConflict 表示尝试创建已存在的资源或违反唯一性约束
var ErrConflict = errors.New("store: item conflict or already exists")

// LinkEntry 代表存储中的一个短链接条目
type LinkEntry struct {
	ShortCode   string    // 短码
	LongURL     string    // 原始长URL (用户输入的，可能需要规范化)
	OriginalURL string    // (可选) 如果对LongURL做了处理，这里存最初的原始输入
	UserID      string    // (可选) 创建用户的ID
	CreatedAt   time.Time // 创建时间
	ExpireAt    time.Time // (可选) 过期时间，零值表示永不过期
	VisitCount  int64     // (可选) 访问计数
}

// Store 定义了数据存储层需要实现的接口
type Store interface {
	// Save 保存短链接映射关系。如果shortCode已存在，应返回ErrConflict。
	Save(ctx context.Context, entry *LinkEntry) error
	// FindByShortCode 根据短码查找原始长链接信息。如果未找到，返回ErrNotFound。
	FindByShortCode(ctx context.Context, shortCode string) (*LinkEntry, error)
	// FindByOriginalURLAndUserID (可选) 根据原始URL和用户ID查找是否已存在短链接 (防止重复创建)
	// FindByOriginalURLAndUserID(ctx context.Context, originalURL string, userID string) (*LinkEntry, error)
	// IncrementVisitCount (可选) 增加短码的访问计数
	// IncrementVisitCount(ctx context.Context, shortCode string) error
	// GetVisitCount (可选) 获取短码的访问计数
	// GetVisitCount(ctx context.Context, shortCode string) (int64, error)

	Close() error // 关闭存储连接或释放资源
}
