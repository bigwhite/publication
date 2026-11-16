package cancellations

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"example.com/go-mcp-enterprise/internal/orders"
)

var (
	ErrInvalidOrderID     = errors.New("order_id is required")
	ErrInvalidReason      = errors.New("reason must contain 1 to 500 characters")
	ErrInvalidIdempotency = errors.New("idempotency_key is required")
)

// OrderReader 是取消申请服务需要的最小订单读取接口。
// 这里不依赖订单服务的具体实现。
type OrderReader interface {
	GetOrder(ctx context.Context, orderID string) (orders.Order, error)
}

type Service struct {
	orderReader OrderReader
	repo        Repository
}

func NewService(
	orderReader OrderReader,
	repo Repository,
) *Service {
	return &Service{
		orderReader: orderReader,
		repo:        repo,
	}
}

// RequestCancellation 创建订单取消申请。
// 它不会直接修改订单状态，也不会触发退款或库存释放。
func (s *Service) RequestCancellation(
	ctx context.Context,
	orderID string,
	reason string,
	idempotencyKey string,
) (*Request, error) {
	orderID = strings.TrimSpace(orderID)
	reason = strings.TrimSpace(reason)
	idempotencyKey = strings.TrimSpace(idempotencyKey)

	if orderID == "" {
		return nil, ErrInvalidOrderID
	}

	if reason == "" || utf8.RuneCountInString(reason) > 500 {
		return nil, ErrInvalidReason
	}

	if idempotencyKey == "" {
		return nil, ErrInvalidIdempotency
	}

	// 先确认订单存在。
	// 是否允许最终取消，应由后续审批和订单领域服务判断。
	if _, err := s.orderReader.GetOrder(ctx, orderID); err != nil {
		return nil, err
	}

	req := &Request{
		OrderID:        orderID,
		Reason:         reason,
		Status:         "pending_approval",
		IdempotencyKey: idempotencyKey,
	}

	return s.repo.Create(ctx, req)
}
