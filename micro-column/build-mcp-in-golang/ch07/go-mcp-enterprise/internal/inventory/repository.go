package inventory

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("inventory item not found")

type Repository interface {
	Get(ctx context.Context, sku string) (*Item, error)
	List(ctx context.Context) ([]Item, error)
}
