package cancellations

import (
	"context"
	"errors"
)

var (
	ErrIdempotencyConflict = errors.New(
		"idempotency key already used with different request",
	)
)

// Repository 定义取消申请的持久化行为。
type Repository interface {
	// Create 创建申请。
	// 如果幂等键已经存在且请求内容一致，则返回已有申请。
	// 如果幂等键已存在但请求内容不同，则返回冲突错误。
	Create(ctx context.Context, req *Request) (*Request, error)
}
