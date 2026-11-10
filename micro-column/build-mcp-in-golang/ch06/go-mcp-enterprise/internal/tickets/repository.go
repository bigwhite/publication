package tickets

import (
	"context"
	"errors"
)

var (
	ErrInvalidInput        = errors.New("invalid ticket input")
	ErrIdempotencyConflict = errors.New("idempotency key conflict")
)

type Repository interface {
	Create(ctx context.Context, input CreateInput) (Ticket, error)
}
