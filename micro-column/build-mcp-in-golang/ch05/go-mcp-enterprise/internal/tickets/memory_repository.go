package tickets

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type createRecord struct {
	input  CreateInput
	ticket Ticket
}

type MemoryRepository struct {
	mu      sync.Mutex
	records map[string]createRecord
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		records: make(map[string]createRecord),
	}
}

func (r *MemoryRepository) Create(
	ctx context.Context,
	input CreateInput,
) (Ticket, error) {
	if err := ctx.Err(); err != nil {
		return Ticket{}, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// 相同幂等键已经执行过。
	if record, ok := r.records[input.IdempotencyKey]; ok {
		// 相同请求：返回原来的工单。
		if record.input == input {
			return record.ticket, nil
		}

		// 相同幂等键却对应不同请求：拒绝。
		return Ticket{}, ErrIdempotencyConflict
	}

	ticket := Ticket{
		ID:          fmt.Sprintf("TKT-%06d", len(r.records)+1),
		OrderID:     input.OrderID,
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		Status:      "open",
		CreatedAt:   time.Now().UTC(),
	}

	r.records[input.IdempotencyKey] = createRecord{
		input:  input,
		ticket: ticket,
	}

	return ticket, nil
}
