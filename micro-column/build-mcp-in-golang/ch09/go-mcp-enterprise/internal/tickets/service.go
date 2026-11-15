package tickets

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"example.com/go-mcp-enterprise/internal/orders"
)

type OrderReader interface {
	GetOrder(ctx context.Context, orderID string) (orders.Order, error)
}

type Service struct {
	orders OrderReader
	repo   Repository
}

func NewService(
	orderReader OrderReader,
	repo Repository,
) *Service {
	return &Service{
		orders: orderReader,
		repo:   repo,
	}
}

func (s *Service) CreateTicket(
	ctx context.Context,
	input CreateInput,
) (Ticket, error) {
	input.OrderID = strings.TrimSpace(input.OrderID)
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)
	input.Priority = strings.TrimSpace(input.Priority)
	input.IdempotencyKey = strings.TrimSpace(input.IdempotencyKey)

	if input.OrderID == "" ||
		input.Title == "" ||
		input.Description == "" ||
		input.IdempotencyKey == "" {
		return Ticket{}, fmt.Errorf("%w: required field is empty", ErrInvalidInput)
	}

	if utf8.RuneCountInString(input.Title) > 120 {
		return Ticket{}, fmt.Errorf("%w: title is too long", ErrInvalidInput)
	}

	if utf8.RuneCountInString(input.Description) > 2000 {
		return Ticket{}, fmt.Errorf("%w: description is too long", ErrInvalidInput)
	}

	switch input.Priority {
	case "low", "normal", "high":
	default:
		return Ticket{}, fmt.Errorf("%w: invalid priority", ErrInvalidInput)
	}

	// 工单必须关联一笔真实存在的订单。
	if _, err := s.orders.GetOrder(ctx, input.OrderID); err != nil {
		return Ticket{}, err
	}

	return s.repo.Create(ctx, input)
}
