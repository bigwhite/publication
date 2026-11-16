package cancellations

import (
	"context"
	"errors"
	"testing"

	"example.com/go-mcp-enterprise/internal/orders"
)

type fakeOrderReader struct {
	err error
}

func (f fakeOrderReader) GetOrder(
	ctx context.Context,
	orderID string,
) (orders.Order, error) {
	if f.err != nil {
		return orders.Order{}, f.err
	}

	// 本测试只验证订单读取是否被调用成功，
	// 不依赖订单对象的具体字段。
	return orders.Order{}, nil
}

func newTestService(reader OrderReader) *Service {
	return NewService(reader, NewMemoryRepository())
}

func TestRequestCancellation(t *testing.T) {
	ctx := context.Background()
	service := newTestService(fakeOrderReader{})

	req, err := service.RequestCancellation(
		ctx,
		"ORD-1001",
		"用户主动提出取消申请",
		"cancel-ORD-1001-001",
	)
	if err != nil {
		t.Fatalf("RequestCancellation() error = %v", err)
	}

	if req.ID == "" {
		t.Fatal("expected non-empty request ID")
	}

	if req.OrderID != "ORD-1001" {
		t.Fatalf("unexpected order ID: %s", req.OrderID)
	}

	if req.Status != "pending_approval" {
		t.Fatalf("unexpected status: %s", req.Status)
	}
}

func TestRequestCancellationIdempotency(t *testing.T) {
	ctx := context.Background()
	service := newTestService(fakeOrderReader{})

	first, err := service.RequestCancellation(
		ctx,
		"ORD-1001",
		"用户主动提出取消申请",
		"cancel-ORD-1001-001",
	)
	if err != nil {
		t.Fatalf("first request: %v", err)
	}

	second, err := service.RequestCancellation(
		ctx,
		"ORD-1001",
		"用户主动提出取消申请",
		"cancel-ORD-1001-001",
	)
	if err != nil {
		t.Fatalf("second request: %v", err)
	}

	if first.ID != second.ID {
		t.Fatalf(
			"expected same request ID, got %q and %q",
			first.ID,
			second.ID,
		)
	}
}

func TestRequestCancellationIdempotencyConflict(t *testing.T) {
	ctx := context.Background()
	service := newTestService(fakeOrderReader{})

	_, err := service.RequestCancellation(
		ctx,
		"ORD-1001",
		"用户主动提出取消申请",
		"cancel-ORD-1001-001",
	)
	if err != nil {
		t.Fatalf("first request: %v", err)
	}

	_, err = service.RequestCancellation(
		ctx,
		"ORD-1001",
		"修改后的取消原因",
		"cancel-ORD-1001-001",
	)
	if !errors.Is(err, ErrIdempotencyConflict) {
		t.Fatalf("expected ErrIdempotencyConflict, got %v", err)
	}
}

func TestRequestCancellationOrderNotFound(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("order not found")

	service := newTestService(fakeOrderReader{
		err: expectedErr,
	})

	_, err := service.RequestCancellation(
		ctx,
		"ORD-9999",
		"用户主动提出取消申请",
		"cancel-ORD-9999-001",
	)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected order error, got %v", err)
	}
}

func TestRequestCancellationValidation(t *testing.T) {
	ctx := context.Background()
	service := newTestService(fakeOrderReader{})

	tests := []struct {
		name           string
		orderID        string
		reason         string
		idempotencyKey string
		wantErr        error
	}{
		{
			name:           "missing order ID",
			reason:         "用户申请取消",
			idempotencyKey: "key-1",
			wantErr:        ErrInvalidOrderID,
		},
		{
			name:           "missing reason",
			orderID:        "ORD-1001",
			idempotencyKey: "key-2",
			wantErr:        ErrInvalidReason,
		},
		{
			name:    "missing idempotency key",
			orderID: "ORD-1001",
			reason:  "用户申请取消",
			wantErr: ErrInvalidIdempotency,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.RequestCancellation(
				ctx,
				tt.orderID,
				tt.reason,
				tt.idempotencyKey,
			)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}
