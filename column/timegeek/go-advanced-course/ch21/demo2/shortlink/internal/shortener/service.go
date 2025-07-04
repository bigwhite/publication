package shortener

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/your_org/shortlink/internal/idgen"
	"github.com/your_org/shortlink/internal/storage"
)

var ErrInvalidLongURL = errors.New("shortener: long URL is invalid or empty")
var ErrShortCodeTooShort = errors.New("shortener: short code is too short or invalid")
var ErrShortCodeGenerationFailed = errors.New("shortener: failed to generate a unique short code after multiple attempts")
var ErrLinkNotFound = errors.New("shortener: link not found")
var ErrConflict = errors.New("shortener: conflict, possibly short code exists or generation failed after retries")

type Config struct {
	Store           storage.Store   // 依赖接口
	Generator       idgen.Generator // 依赖接口
	Logger          *log.Logger     // 接收标准库 logger
	MaxGenAttempts  int
	MinShortCodeLen int
}

type Service struct {
	store           storage.Store
	generator       idgen.Generator
	logger          *log.Logger
	maxGenAttempts  int
	minShortCodeLen int
}

func NewService(cfg Config) (*Service, error) {
	if cfg.Store == nil {
		return nil, errors.New("shortener: store is required for service")
	}
	if cfg.Generator == nil {
		return nil, errors.New("shortener: generator is required for service")
	}
	if cfg.Logger == nil {
		cfg.Logger = log.New(os.Stdout, "[ShortenerService-Default] ", log.LstdFlags|log.Lshortfile)
	}
	if cfg.MaxGenAttempts <= 0 {
		cfg.MaxGenAttempts = 3
	}
	if cfg.MinShortCodeLen <= 0 {
		cfg.MinShortCodeLen = 5
	}
	return &Service{
		store:           cfg.Store,
		generator:       cfg.Generator,
		logger:          cfg.Logger,
		maxGenAttempts:  cfg.MaxGenAttempts,
		minShortCodeLen: cfg.MinShortCodeLen,
	}, nil
}

func (s *Service) CreateShortLink(ctx context.Context, longURL string) (string, error) {
	if strings.TrimSpace(longURL) == "" {
		return "", ErrInvalidLongURL
	}

	var shortCode string

	for attempt := 0; attempt < s.maxGenAttempts; attempt++ {
		s.logger.Printf("DEBUG: Attempting to generate short code, attempt %d, longURL_preview: %s\n", attempt+1, preview(longURL, 50))

		code, genErr := s.generator.GenerateShortCode(ctx, longURL)
		if genErr != nil {
			return "", fmt.Errorf("attempt %d to generate short code failed: %w", attempt+1, genErr)
		}
		shortCode = code

		if len(shortCode) < s.minShortCodeLen {
			s.logger.Printf("WARN: Generated short code too short, retrying. Code: %s, Attempt: %d\n", shortCode, attempt+1)
			if attempt < s.maxGenAttempts-1 {
				continue
			} else {
				break
			}
		}

		linkToSave := storage.Link{
			ShortCode: shortCode,
			LongURL:   longURL,
			CreatedAt: time.Now().UTC(),
		}
		saveErr := s.store.Save(ctx, linkToSave)
		if saveErr == nil {
			s.logger.Printf("INFO: Successfully created short link. ShortCode: %s, LongURL_preview: %s\n", shortCode, preview(longURL, 50))
			return shortCode, nil
		}

		if errors.Is(saveErr, storage.ErrShortCodeExists) && attempt < s.maxGenAttempts-1 {
			s.logger.Printf("WARN: Short code collision, retrying. Code: %s, Attempt: %d\n", shortCode, attempt+1)
			continue
		}
		s.logger.Printf("ERROR: Failed to save link. LongURL_preview: %s, ShortCode: %s, Attempt: %d, Error: %v\n", preview(longURL, 50), shortCode, attempt+1, saveErr)
		return "", fmt.Errorf("%w: after %d attempts for input: %w", ErrShortCodeGenerationFailed, attempt+1, saveErr)
	}
	return "", ErrShortCodeGenerationFailed
}

func (s *Service) GetAndTrackLongURL(ctx context.Context, shortCode string) (string, error) {
	if len(shortCode) < s.minShortCodeLen {
		return "", ErrShortCodeTooShort
	}

	link, findErr := s.store.FindByShortCode(ctx, shortCode)
	if findErr != nil {
		if errors.Is(findErr, storage.ErrNotFound) {
			s.logger.Printf("INFO: Short code not found in store. ShortCode: %s\n", shortCode)
			return "", fmt.Errorf("for code '%s': %w", shortCode, ErrLinkNotFound)
		}
		s.logger.Printf("ERROR: Failed to find link by short code. ShortCode: %s, Error: %v\n", shortCode, findErr)
		return "", fmt.Errorf("failed to find link for code '%s': %w", shortCode, findErr)
	}

	go func(sc string, currentCount int64, parentLogger *log.Logger) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		logger := parentLogger // 使用传入的logger实例

		if err := s.store.IncrementVisitCount(bgCtx, sc); err != nil {
			logger.Printf("ERROR: Failed to increment visit count (async). ShortCode: %s, Error: %v\n", sc, err)
		} else {
			logger.Printf("DEBUG: Incremented visit count (async). ShortCode: %s, NewCount_approx: %d\n", sc, currentCount+1)
		}
	}(shortCode, link.VisitCount, s.logger) // 将logger传递进去

	s.logger.Printf("INFO: Redirecting to long URL. ShortCode: %s, LongURL_preview: %s\n", shortCode, preview(link.LongURL, 50))
	return link.LongURL, nil
}

func preview(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}
