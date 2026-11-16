package mcpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"example.com/go-mcp-enterprise/internal/observability"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ObservedGetOrderInput struct {
	OrderID string `json:"order_id" jsonschema:"订单ID"`
}

type ObservedGetOrderOutput struct {
	OrderID string `json:"order_id" jsonschema:"订单ID"`
	Found   bool   `json:"found" jsonschema:"是否找到订单"`
}

// RegisterObservedOrderTool 注册一个用于演示可观测性的订单查询 Tool。
//
// lookup 由调用方提供，避免本文件依赖具体的订单模型。
// 实际项目中可以把同样的 WrapTool 用在已有 get_order Handler 上。
func RegisterObservedOrderTool(
	server *mcp.Server,
	lookup func(context.Context, string) error,
	logger *slog.Logger,
	metrics *observability.Metrics,
) {

	if metrics != nil {
		metrics.RegisterTool("get_order_observed")
	}

	handler := func(
		ctx context.Context,
		_ *mcp.CallToolRequest,
		input ObservedGetOrderInput,
	) (*mcp.CallToolResult, ObservedGetOrderOutput, error) {
		orderID := strings.TrimSpace(input.OrderID)
		if orderID == "" {
			return nil, ObservedGetOrderOutput{}, errors.New("order_id is required")
		}

		if lookup == nil {
			return nil, ObservedGetOrderOutput{}, errors.New("order lookup is not configured")
		}

		if err := lookup(ctx, orderID); err != nil {
			return nil, ObservedGetOrderOutput{},
				fmt.Errorf("get order %q: %w", orderID, err)
		}

		return nil, ObservedGetOrderOutput{
			OrderID: orderID,
			Found:   true,
		}, nil
	}

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "get_order_observed",
			Description: "查询订单，并记录 Tool 调用日志和运行指标",
		},
		observability.WrapTool(
			logger,
			metrics,
			"get_order_observed",
			handler,
		),
	)
}
