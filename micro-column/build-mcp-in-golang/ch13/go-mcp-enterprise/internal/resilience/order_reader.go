package resilience

import (
	"context"
	"time"

	"example.com/go-mcp-enterprise/internal/orders"
)

type OrderReader interface {
	GetOrder(
		ctx context.Context,
		orderID string,
	) (orders.Order, error)
}

// DelayedOrderReader 为本地演示增加可取消的延迟。
// 生产环境不应使用人为延迟包装真实业务服务。
type DelayedOrderReader struct {
	next  OrderReader
	delay time.Duration
}

func NewDelayedOrderReader(
	next OrderReader,
	delay time.Duration,
) *DelayedOrderReader {
	return &DelayedOrderReader{
		next:  next,
		delay: delay,
	}
}

func (r *DelayedOrderReader) GetOrder(
	ctx context.Context,
	orderID string,
) (orders.Order, error) {
	if r.delay > 0 {
		timer := time.NewTimer(r.delay)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			return orders.Order{}, ctx.Err()
		case <-timer.C:
		}
	}

	return r.next.GetOrder(ctx, orderID)
}
