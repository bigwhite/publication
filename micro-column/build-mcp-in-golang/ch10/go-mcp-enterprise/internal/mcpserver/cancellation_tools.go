package mcpserver

import (
	"context"

	"example.com/go-mcp-enterprise/internal/cancellations"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type requestOrderCancellationInput struct {
	OrderID        string `json:"order_id" jsonschema:"需要申请取消的订单 ID"`
	Reason         string `json:"reason" jsonschema:"取消申请原因"`
	IdempotencyKey string `json:"idempotency_key" jsonschema:"本次逻辑操作的幂等键，重试时必须保持不变"`
}

type requestOrderCancellationOutput struct {
	RequestID string `json:"request_id"`
	OrderID   string `json:"order_id"`
	Status    string `json:"status"`
}

// RegisterCancellationTools 注册取消申请 Tool。
// 该 Tool 只创建待审批申请，不会直接取消订单。
func RegisterCancellationTools(
	server *mcp.Server,
	service *cancellations.Service,
) {
	destructive := false
	openWorld := false

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:  "request_order_cancellation",
			Title: "提交订单取消申请",
			Description: "为指定订单提交取消申请。该操作只创建一条待审批申请，" +
				"不会立即取消订单，也不会触发退款、库存释放或物流操作。" +
				"调用前必须确认用户明确提出了取消申请，并提供订单 ID 和取消原因。",
			Annotations: &mcp.ToolAnnotations{
				ReadOnlyHint:    false,
				DestructiveHint: &destructive,
				IdempotentHint:  true,
				OpenWorldHint:   &openWorld,
			},
		},
		func(
			ctx context.Context,
			_ *mcp.CallToolRequest,
			in requestOrderCancellationInput,
		) (*mcp.CallToolResult, requestOrderCancellationOutput, error) {
			req, err := service.RequestCancellation(
				ctx,
				in.OrderID,
				in.Reason,
				in.IdempotencyKey,
			)
			if err != nil {
				return nil, requestOrderCancellationOutput{}, err
			}

			return nil, requestOrderCancellationOutput{
				RequestID: req.ID,
				OrderID:   req.OrderID,
				Status:    req.Status,
			}, nil
		},
	)
}
