package orders

import (
	"context"
	"strings"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetOrder(
	ctx context.Context,
	orderID string,
) (Order, error) {
	orderID = strings.TrimSpace(orderID)
	if orderID == "" {
		return Order{}, ErrInvalidOrderID
	}

	return s.repo.GetByID(ctx, orderID)
}
