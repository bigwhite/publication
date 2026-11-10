package orders

import (
	"context"
	"sync"
	"time"
)

type MemoryRepository struct {
	mu     sync.RWMutex
	orders map[string]Order
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		orders: map[string]Order{
			"ORD-1001": {
				ID:          "ORD-1001",
				CustomerID:  "CUST-001",
				Status:      "paid",
				TotalAmount: 12990,
				Currency:    "CNY",
				CreatedAt:   time.Date(2026, 11, 1, 10, 30, 0, 0, time.UTC),
			},
			"ORD-1002": {
				ID:          "ORD-1002",
				CustomerID:  "CUST-002",
				Status:      "shipped",
				TotalAmount: 25800,
				Currency:    "CNY",
				CreatedAt:   time.Date(2026, 11, 2, 8, 15, 0, 0, time.UTC),
			},
		},
	}
}

func (r *MemoryRepository) GetByID(
	ctx context.Context,
	id string,
) (Order, error) {
	if err := ctx.Err(); err != nil {
		return Order{}, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[id]
	if !ok {
		return Order{}, ErrNotFound
	}

	return order, nil
}
