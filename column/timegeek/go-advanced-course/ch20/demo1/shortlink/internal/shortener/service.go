package shortener

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/your_org/shortlink/internal/idgen"
	"github.com/your_org/shortlink/internal/storage"
)

var ErrInvalidInput = errors.New("shortener: invalid input")
var ErrServiceInternal = errors.New("shortener: internal service error")
var ErrConflict = errors.New("shortener: conflict, possibly short code exists or generation failed after retries")

// Config for Service, allowing dependencies to be passed in.
type Config struct {
	Store          *storage.Store   // Concrete type for demo1
	Generator      *idgen.Generator // Concrete type for demo1
	MaxGenAttempts int
}

type Service struct {
	store          *storage.Store
	generator      *idgen.Generator
	maxGenAttempts int
}

func NewService(cfg Config) *Service { // Removed error return for simplicity in demo1 main
	if cfg.Store == nil || cfg.Generator == nil {
		log.Fatalln("Store and Generator must not be nil for Shortener Service")
	}
	if cfg.MaxGenAttempts <= 0 {
		cfg.MaxGenAttempts = 3
	}
	return &Service{
		store:          cfg.Store,
		generator:      cfg.Generator,
		maxGenAttempts: cfg.MaxGenAttempts,
	}
}

func (s *Service) CreateShortLink(ctx context.Context, longURL string) (string, error) {
	if longURL == "" {
		return "", fmt.Errorf("%w: longURL cannot be empty", ErrInvalidInput)
	}

	var shortCode string
	for attempt := 0; attempt < s.maxGenAttempts; attempt++ {
		log.Printf("DEBUG: Attempting to generate short code, attempt %d, longURL: %s\n", attempt+1, longURL)
		code, genErr := s.generator.GenerateShortCode(ctx, longURL)
		if genErr != nil {
			return "", fmt.Errorf("attempt %d: failed to generate short code: %w", attempt+1, genErr)
		}
		shortCode = code

		linkToSave := storage.Link{
			ShortCode: shortCode,
			LongURL:   longURL,
			CreatedAt: time.Now().UTC(),
		}
		saveErr := s.store.Save(ctx, linkToSave)
		if saveErr == nil {
			log.Printf("INFO: Successfully created short link. ShortCode: %s, LongURL: %s\n", shortCode, longURL)
			return shortCode, nil
		}

		if errors.Is(saveErr, storage.ErrShortCodeExists) && attempt < s.maxGenAttempts-1 {
			log.Printf("WARN: Short code collision, retrying. Code: %s, Attempt: %d\n", shortCode, attempt+1)
			continue
		}
		log.Printf("ERROR: Failed to save link. LongURL: %s, ShortCode: %s, Attempt: %d, Error: %v\n", longURL, shortCode, attempt+1, saveErr)
		return "", fmt.Errorf("%w: after %d attempts for input %s: %w", ErrConflict, attempt+1, longURL, saveErr)
	}
	return "", ErrServiceInternal
}

func (s *Service) GetAndTrackLongURL(ctx context.Context, shortCode string) (string, error) {
	if shortCode == "" {
		return "", fmt.Errorf("%w: short code cannot be empty", ErrInvalidInput)
	}

	link, findErr := s.store.FindByShortCode(ctx, shortCode)
	if findErr != nil {
		if errors.Is(findErr, storage.ErrNotFound) {
			log.Printf("WARN: Short code not found. ShortCode: %s\n", shortCode)
			return "", fmt.Errorf("short link for code '%s' not found: %w", shortCode, findErr)
		}
		log.Printf("ERROR: Failed to find link by short code. ShortCode: %s, Error: %v\n", shortCode, findErr)
		return "", fmt.Errorf("failed to find link for code '%s': %w", shortCode, findErr)
	}

	go func(sc string, currentCount int64) {
		bgCtx := context.Background()
		if err := s.store.IncrementVisitCount(bgCtx, sc); err != nil {
			log.Printf("ERROR: Failed to increment visit count (async). ShortCode: %s, Error: %v\n", sc, err)
		} else {
			log.Printf("DEBUG: Incremented visit count (async). ShortCode: %s, NewCount: %d\n", sc, currentCount+1)
		}
	}(shortCode, link.VisitCount)

	log.Printf("INFO: Redirecting to long URL. ShortCode: %s, LongURL: %s\n", shortCode, link.LongURL)
	return link.LongURL, nil
}
