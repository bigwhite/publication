package mcpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"example.com/go-mcp-enterprise/internal/orders"
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

// New 创建企业订单系统的 MCP Server。
func New(service OrderService, logger *slog.Logger) *mcp.Server {
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
			Description: "根据订单编号查询订单详细信息。该工具只读取订单，不会修改订单状态。",
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

	return server
}
