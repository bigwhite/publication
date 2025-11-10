// pkg/repository/postgres/link_repository.go
package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bigwhite/shortlink/pkg/domain"
	"github.com/bigwhite/shortlink/pkg/repository"
	"github.com/lib/pq"
)

type PgLinkRepository struct {
	db repository.DBTX
}

func NewPgLinkRepository(db repository.DBTX) *PgLinkRepository {
	return &PgLinkRepository{db: db}
}

func (p *PgLinkRepository) FindByCode(ctx context.Context, code string) (*domain.Link, error) {
	var link domain.Link
	query := "SELECT id, original_url, short_code, created_at FROM links WHERE short_code=$1"

	// 使用 QueryRowContext，然后手动 Scan
	row := p.db.QueryRowContext(ctx, query, code)
	err := row.Scan(&link.ID, &link.OriginalURL, &link.ShortCode, &link.CreatedAt) // <-- 关键变更：手动按顺序 Scan

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // 找不到不是错误
		}
		return nil, err // 其他数据库错误
	}
	return &link, nil
}

func (p *PgLinkRepository) Save(ctx context.Context, link *domain.Link) error {
	query := "INSERT INTO links (original_url, short_code) VALUES ($1, $2) RETURNING id, created_at"

	row := p.db.QueryRowContext(ctx, query, link.OriginalURL, link.ShortCode)
	err := row.Scan(&link.ID, &link.CreatedAt) // <-- 关键变更：手动 Scan 返回的字段

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return repository.ErrCodeConflict
		}
		return err
	}
	return nil
}
