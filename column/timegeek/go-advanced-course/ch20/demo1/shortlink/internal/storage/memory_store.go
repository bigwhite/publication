package storage

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

type Link struct {
	ShortCode  string
	LongURL    string
	VisitCount int64
	CreatedAt  time.Time
}

var ErrNotFound = errors.New("storage: link not found")
var ErrShortCodeExists = errors.New("storage: short code already exists")

type Store struct {
	mu    sync.RWMutex
	links map[string]Link
}

func NewStore() *Store {
	return &Store{
		links: make(map[string]Link),
	}
}

func (s *Store) Save(ctx context.Context, link Link) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.links[link.ShortCode]; exists {
		return ErrShortCodeExists
	}
	if link.CreatedAt.IsZero() {
		link.CreatedAt = time.Now().UTC()
	}
	s.links[link.ShortCode] = link
	log.Printf("DEBUG: Saved link. ShortCode: %s, LongURL: %s\n", link.ShortCode, link.LongURL)
	return nil
}

func (s *Store) FindByShortCode(ctx context.Context, shortCode string) (*Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	link, exists := s.links[shortCode]
	if !exists {
		return nil, ErrNotFound
	}
	linkCopy := link
	log.Printf("DEBUG: Found link. ShortCode: %s, LongURL: %s\n", link.ShortCode, link.LongURL)
	return &linkCopy, nil
}

func (s *Store) IncrementVisitCount(ctx context.Context, shortCode string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	link, exists := s.links[shortCode]
	if !exists {
		return ErrNotFound
	}
	link.VisitCount++
	s.links[shortCode] = link
	log.Printf("DEBUG: Incremented visit count. ShortCode: %s, NewCount: %d\n", link.ShortCode, link.VisitCount)
	return nil
}

func (s *Store) Close() error {
	log.Println("MemoryStore closed (no-op)")
	return nil
}
