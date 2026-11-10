package orders

import (
	"context"
	"errors"
)

var (
	ErrNotFound       = errors.New("order not found")
	ErrInvalidOrderID = errors.New("invalid order id")
)

type Repository interface {
	GetByID(ctx context.Context, id string) (Order, error)
}
