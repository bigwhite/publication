package orders

import (
	"context"
	"errors"
	"testing"
)

func TestServiceGetOrder(t *testing.T) {
	service := NewService(NewMemoryRepository())

	order, err := service.GetOrder(
		context.Background(),
		"ORD-1001",
	)
	if err != nil {
		t.Fatalf("GetOrder() error = %v", err)
	}

	if order.ID != "ORD-1001" {
		t.Fatalf("unexpected order ID: %s", order.ID)
	}

	if order.Status != "paid" {
		t.Fatalf("unexpected order status: %s", order.Status)
	}

	_, err = service.GetOrder(
		context.Background(),
		"ORD-404",
	)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	_, err = service.GetOrder(
		context.Background(),
		"",
	)
	if !errors.Is(err, ErrInvalidOrderID) {
		t.Fatalf("expected ErrInvalidOrderID, got %v", err)
	}
}
