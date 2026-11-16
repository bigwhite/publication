package cancellations

import (
	"context"
	"fmt"
	"sync"
)

type MemoryRepository struct {
	mu sync.Mutex

	byKey map[string]*Request
	next  uint64
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		byKey: make(map[string]*Request),
	}
}

func (r *MemoryRepository) Create(
	ctx context.Context,
	req *Request,
) (*Request, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// 已经存在相同幂等键时，检查请求内容。
	if existing, ok := r.byKey[req.IdempotencyKey]; ok {
		if existing.OrderID != req.OrderID ||
			existing.Reason != req.Reason {
			return nil, ErrIdempotencyConflict
		}

		// 相同请求返回之前创建的申请。
		return cloneRequest(existing), nil
	}

	r.next++

	created := &Request{
		ID:             fmt.Sprintf("CRQ-%06d", r.next),
		OrderID:        req.OrderID,
		Reason:         req.Reason,
		Status:         "pending_approval",
		IdempotencyKey: req.IdempotencyKey,
	}

	r.byKey[created.IdempotencyKey] = created

	return cloneRequest(created), nil
}

func cloneRequest(req *Request) *Request {
	if req == nil {
		return nil
	}

	copied := *req
	return &copied
}
