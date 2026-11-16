package inventory

import (
	"context"
	"errors"
	"testing"
)

func TestGetItem(t *testing.T) {
	service := NewService(NewMemoryRepository())

	item, err := service.GetItem(
		context.Background(),
		"SKU-1002",
	)
	if err != nil {
		t.Fatalf("GetItem() error = %v", err)
	}

	if item.SKU != "SKU-1002" {
		t.Fatalf("SKU = %q, want SKU-1002", item.SKU)
	}

	if item.AvailableQuantity != 8 {
		t.Fatalf(
			"AvailableQuantity = %d, want 8",
			item.AvailableQuantity,
		)
	}
}

func TestGetItemNotFound(t *testing.T) {
	service := NewService(NewMemoryRepository())

	_, err := service.GetItem(
		context.Background(),
		"SKU-UNKNOWN",
	)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("error = %v, want ErrNotFound", err)
	}
}

func TestGetItemInvalidSKU(t *testing.T) {
	service := NewService(NewMemoryRepository())

	_, err := service.GetItem(
		context.Background(),
		"   ",
	)
	if !errors.Is(err, ErrInvalidSKU) {
		t.Fatalf("error = %v, want ErrInvalidSKU", err)
	}
}

func TestListLowStock(t *testing.T) {
	service := NewService(NewMemoryRepository())

	items, err := service.ListLowStock(
		context.Background(),
		10,
	)
	if err != nil {
		t.Fatalf("ListLowStock() error = %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("len(items) = %d, want 2", len(items))
	}

	if items[0].SKU != "SKU-1002" {
		t.Fatalf(
			"items[0].SKU = %q, want SKU-1002",
			items[0].SKU,
		)
	}

	if items[1].SKU != "SKU-1003" {
		t.Fatalf(
			"items[1].SKU = %q, want SKU-1003",
			items[1].SKU,
		)
	}
}

func TestListLowStockInvalidThreshold(t *testing.T) {
	service := NewService(NewMemoryRepository())

	_, err := service.ListLowStock(
		context.Background(),
		-1,
	)
	if !errors.Is(err, ErrInvalidThreshold) {
		t.Fatalf(
			"error = %v, want ErrInvalidThreshold",
			err,
		)
	}
}
