// pkg/repository/link.go
package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bigwhite/shortlink/pkg/domain"
)

var ErrCodeConflict = errors.New("short code conflicts")

// DBTX 是一个接口，抽象了 *sql.DB 和 *sql.Tx 的共同方法
// 它只包含我们 Repository 需要用到的方法
type DBTX interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type LinkRepository interface {
	FindByCode(ctx context.Context, code string) (*domain.Link, error)
	Save(ctx context.Context, link *domain.Link) error
}

// LinkCache 定义了链接统计的缓存接口
type LinkCache interface {
	IncrementVisitCount(ctx context.Context, code string) error
	GetVisitCount(ctx context.Context, code string) (int64, error)
}
