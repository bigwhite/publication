// pkg/service/shortener.go
package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"time"

	"github.com/bigwhite/shortlink/pkg/domain"
	"github.com/bigwhite/shortlink/pkg/repository"
	"github.com/bigwhite/shortlink/pkg/validator"
)

// generateShortCode 是生成短码的默认实现
func generateShortCode(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// ShortenerService 提供了短链接的核心业务逻辑
type ShortenerService struct {
	repo                  repository.LinkRepository
	cache                 repository.LinkCache
	generateShortCodeFunc func(int) (string, error) // 短码生成函数的依赖
}

// NewShortenerService 是 ShortenerService 的构造函数
// 它为 generateShortCodeFunc 字段设置了默认的实现
func NewShortenerService(repo repository.LinkRepository, cache repository.LinkCache) *ShortenerService {
	return &ShortenerService{
		repo:                  repo,
		cache:                 cache,
		generateShortCodeFunc: generateShortCode,
	}
}

// CreateLink 创建一个新的短链接
func (s *ShortenerService) CreateLink(ctx context.Context, originalURL string) (*domain.Link, error) {
	// 1. 验证 URL 的合法性
	if !validator.IsValidURL(originalURL) {
		return nil, errors.New("invalid URL")
	}

	// 为数据库操作创建一个带超时的 context
	// 这是一个很好的实践，防止下游依赖的缓慢拖垮整个服务
	dbCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// 2. 尝试生成一个唯一的短码
	const maxRetries = 5
	for i := 0; i < maxRetries; i++ {
		// 调用注入的函数来生成短码，而不是包级别的函数
		code, err := s.generateShortCodeFunc(6)
		if err != nil {
			return nil, err
		}

		link := &domain.Link{
			OriginalURL: originalURL,
			ShortCode:   code,
		}

		// 3. 尝试保存到仓库，使用带超时的 context
		err = s.repo.Save(dbCtx, link) // <-- 使用 dbCtx
		if err == nil {
			return link, nil // 保存成功，直接返回
		}

		// 4. 如果不是冲突错误，则直接返回错误
		if !errors.Is(err, repository.ErrCodeConflict) {
			return nil, err
		}

		// 5. 如果是冲突错误，则继续下一次循环
		// (可以在这里增加日志，记录发生了一次冲突)
	}

	// 6. 如果重试次数耗尽，返回失败错误
	return nil, errors.New("failed to create a unique short code after multiple retries")
}

// Redirect 处理重定向逻辑
func (s *ShortenerService) Redirect(ctx context.Context, code string) (*domain.Link, error) {
	link, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if link == nil {
		return nil, errors.New("link not found")
	}

	// 异步地增加访问计数 (使用 goroutine)
	go func() {
		err := s.cache.IncrementVisitCount(context.Background(), code)
		if err != nil {
			// 在真实应用中，这里应该使用结构化日志记录错误
			log.Printf("ERROR: Failed to increment visit count for code %s: %v", code, err)
		}
	}()

	return link, nil
}


func (s *ShortenerService) GetStats(ctx context.Context, code string) (int64, error) {
	// 直接从缓存中读取
	return s.cache.GetVisitCount(ctx, code)
}
