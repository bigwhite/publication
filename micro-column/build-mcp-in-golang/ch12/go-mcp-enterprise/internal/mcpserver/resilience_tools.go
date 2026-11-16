package mcpserver

import (
	"context"
	"errors"
	"time"

	"example.com/go-mcp-enterprise/internal/orders"
	"example.com/go-mcp-enterprise/internal/resilience"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ResilientGetOrderInput struct {
	OrderID string `json:"order_id" jsonschema:"订单 ID"`
}

type ResilientGetOrderOutput struct {
	Order orders.Order `json:"order"`
}

func RegisterResilientOrderTools(
	server *mcp.Server,
	reader resilience.OrderReader,
	gate *resilience.Gate,
	limiter *resilience.Limiter,
) {
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name: "get_order_resilient",
			Description: "在并发、速率和执行时间限制下查询订单。" +
				"仅查询当前业务服务允许访问的订单。",
		},
		func(
			ctx context.Context,
			req *mcp.CallToolRequest,
			input ResilientGetOrderInput,
		) (*mcp.CallToolResult, ResilientGetOrderOutput, error) {
			if input.OrderID == "" {
				return nil, ResilientGetOrderOutput{}, errors.New(
					"order_id is required",
				)
			}

			if !limiter.Allow() {
				return nil, ResilientGetOrderOutput{}, errors.New(
					"rate limit exceeded",
				)
			}

			if err := gate.Acquire(ctx); err != nil {
				return nil, ResilientGetOrderOutput{}, err
			}
			defer gate.Release()

			callCtx, cancel := context.WithTimeout(
				ctx,
				2*time.Second,
			)
			defer cancel()

			order, err := reader.GetOrder(callCtx, input.OrderID)
			if err != nil {
				return nil, ResilientGetOrderOutput{}, err
			}

			return nil, ResilientGetOrderOutput{
				Order: order,
			}, nil
		},
	)
}
