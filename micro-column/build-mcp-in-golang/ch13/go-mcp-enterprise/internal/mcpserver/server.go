package mcpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"example.com/go-mcp-enterprise/internal/orders"
	"example.com/go-mcp-enterprise/internal/tickets"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type OrderService interface {
	GetOrder(ctx context.Context, orderID string) (orders.Order, error)
}

type GetOrderInput struct {
	OrderID string `json:"order_id" jsonschema:"the unique order identifier, such as ORD-1001"`
}

type GetOrderOutput struct {
	ID          string `json:"id" jsonschema:"the order identifier"`
	CustomerID  string `json:"customer_id" jsonschema:"the customer identifier"`
	Status      string `json:"status" jsonschema:"the current order status"`
	TotalAmount int64  `json:"total_amount" jsonschema:"the amount in the smallest currency unit"`
	Currency    string `json:"currency" jsonschema:"the currency code"`
	CreatedAt   string `json:"created_at" jsonschema:"the order creation time in RFC3339 format"`
}

type CreateTicketInput struct {
	OrderID        string `json:"order_id" jsonschema:"existing order identifier"`
	Title          string `json:"title" jsonschema:"ticket title, maximum 120 characters"`
	Description    string `json:"description" jsonschema:"problem description, maximum 2000 characters"`
	Priority       string `json:"priority" jsonschema:"priority: low, normal, or high"`
	IdempotencyKey string `json:"idempotency_key" jsonschema:"stable key for this logical request; reuse it when retrying"`
}

type CreateTicketOutput struct {
	TicketID  string `json:"ticket_id" jsonschema:"created ticket identifier"`
	OrderID   string `json:"order_id" jsonschema:"associated order identifier"`
	Status    string `json:"status" jsonschema:"ticket status"`
	Priority  string `json:"priority" jsonschema:"ticket priority"`
	CreatedAt string `json:"created_at" jsonschema:"ticket creation time in RFC3339 format"`
}

type TicketService interface {
	CreateTicket(ctx context.Context, input tickets.CreateInput) (tickets.Ticket, error)
}

func boolPtr(v bool) *bool {
	return &v
}

// New 创建企业订单系统的 MCP Server。
func New(
	service OrderService,
	ticketService TicketService,
	logger *slog.Logger,
) *mcp.Server {
	if logger == nil {
		logger = slog.Default()
	}

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "go-mcp-enterprise",
			Version: "0.1.0",
		},
		&mcp.ServerOptions{
			Instructions: "企业订单系统 MCP 服务，提供订单查询能力。",
			Logger:       logger,
			Capabilities: &mcp.ServerCapabilities{},
		},
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "get_order",
			Title:       "查询订单",
			Description: "根据唯一订单编号查询订单详细信息。只读取数据，不修改订单。",
			Annotations: &mcp.ToolAnnotations{
				ReadOnlyHint:  true,
				OpenWorldHint: boolPtr(false),
			},
		},
		func(
			ctx context.Context,
			req *mcp.CallToolRequest,
			input GetOrderInput,
		) (*mcp.CallToolResult, GetOrderOutput, error) {
			order, err := service.GetOrder(ctx, input.OrderID)
			if err != nil {
				switch {
				case errors.Is(err, orders.ErrInvalidOrderID):
					return nil, GetOrderOutput{},
						errors.New("订单编号不能为空")

				case errors.Is(err, orders.ErrNotFound):
					return nil, GetOrderOutput{},
						fmt.Errorf("订单 %q 不存在", input.OrderID)

				case errors.Is(err, context.Canceled),
					errors.Is(err, context.DeadlineExceeded):
					return nil, GetOrderOutput{}, err

				default:
					logger.ErrorContext(
						ctx,
						"get order failed",
						"order_id", input.OrderID,
						"error", err,
					)

					return nil, GetOrderOutput{},
						errors.New("订单查询暂时不可用，请稍后重试")
				}
			}

			output := GetOrderOutput{
				ID:          order.ID,
				CustomerID:  order.CustomerID,
				Status:      order.Status,
				TotalAmount: order.TotalAmount,
				Currency:    order.Currency,
				CreatedAt:   order.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
			}

			return nil, output, nil
		},
	)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_ticket",
		Title:       "创建售后工单",
		Description: "为已存在的订单创建一张售后支持工单。需要提供订单编号、问题标题、问题描述和优先级。该操作会创建新工单，但不会取消订单、退款或修改订单状态。",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    false,
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(false),
		},
	}, func(
		ctx context.Context,
		req *mcp.CallToolRequest,
		input CreateTicketInput,
	) (*mcp.CallToolResult, CreateTicketOutput, error) {
		ticket, err := ticketService.CreateTicket(
			ctx,
			tickets.CreateInput{
				OrderID:        input.OrderID,
				Title:          input.Title,
				Description:    input.Description,
				Priority:       input.Priority,
				IdempotencyKey: input.IdempotencyKey,
			},
		)
		if err != nil {
			switch {
			case errors.Is(err, tickets.ErrInvalidInput):
				return nil, CreateTicketOutput{}, err

			case errors.Is(err, tickets.ErrIdempotencyConflict):
				return nil, CreateTicketOutput{},
					errors.New("幂等键已用于不同的工单请求")

			case errors.Is(err, orders.ErrNotFound):
				return nil, CreateTicketOutput{},
					errors.New("关联订单不存在")

			default:
				logger.ErrorContext(
					ctx,
					"create ticket failed",
					"order_id", input.OrderID,
					"error", err,
				)
				return nil, CreateTicketOutput{},
					errors.New("工单创建暂时不可用，请稍后重试")
			}
		}

		return nil, CreateTicketOutput{
			TicketID:  ticket.ID,
			OrderID:   ticket.OrderID,
			Status:    ticket.Status,
			Priority:  ticket.Priority,
			CreatedAt: ticket.CreatedAt.Format(time.RFC3339),
		}, nil
	})

	return server
}
