package inventory

import (
	"context"
	"sort"
	"sync"
)

type MemoryRepository struct {
	mu    sync.RWMutex
	items map[string]Item
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		items: map[string]Item{
			"SKU-1001": {
				SKU:               "SKU-1001",
				Name:              "Mechanical Keyboard",
				TotalQuantity:     100,
				AvailableQuantity: 68,
				ReservedQuantity:  32,
				ReorderPoint:      20,
			},
			"SKU-1002": {
				SKU:               "SKU-1002",
				Name:              "Wireless Mouse",
				TotalQuantity:     80,
				AvailableQuantity: 8,
				ReservedQuantity:  72,
				ReorderPoint:      15,
			},
			"SKU-1003": {
				SKU:               "SKU-1003",
				Name:              "4K Monitor",
				TotalQuantity:     30,
				AvailableQuantity: 3,
				ReservedQuantity:  27,
				ReorderPoint:      5,
			},
		},
	}
}

func (r *MemoryRepository) Get(
	ctx context.Context,
	sku string,
) (*Item, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.items[sku]
	if !ok {
		return nil, ErrNotFound
	}

	result := item
	return &result, nil
}

func (r *MemoryRepository) List(
	ctx context.Context,
) ([]Item, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Item, 0, len(r.items))
	for _, item := range r.items {
		result = append(result, item)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].SKU < result[j].SKU
	})

	return result, nil
}
