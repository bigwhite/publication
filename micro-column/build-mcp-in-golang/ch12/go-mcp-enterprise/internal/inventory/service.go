package inventory

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrInvalidSKU       = errors.New("invalid sku")
	ErrInvalidThreshold = errors.New("invalid threshold")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetItem(
	ctx context.Context,
	sku string,
) (*Item, error) {
	sku = strings.TrimSpace(sku)
	if sku == "" {
		return nil, ErrInvalidSKU
	}

	return s.repo.Get(ctx, sku)
}

func (s *Service) ListLowStock(
	ctx context.Context,
	threshold int,
) ([]Item, error) {
	if threshold < 0 {
		return nil, ErrInvalidThreshold
	}

	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]Item, 0)
	for _, item := range items {
		if item.AvailableQuantity <= threshold {
			result = append(result, item)
		}
	}

	return result, nil
}
