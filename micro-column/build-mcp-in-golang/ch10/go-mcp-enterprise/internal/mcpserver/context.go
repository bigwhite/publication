package mcpserver

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	orderStatusURI = "biz://docs/order-status"
)

const orderStatusGuide = `# 订单状态说明

订单状态及其含义：

- pending：订单已创建，尚未完成支付。
- paid：订单已支付，等待后续履约。
- shipped：订单已发货，正在配送途中。
- completed：订单已完成。
- cancelled：订单已取消。

使用规则：

1. 本文档解释订单状态的业务含义，不代表某个订单的实时状态。
2. 订单当前状态必须以 get_order Tool 返回的数据为准。
3. 不要根据用户描述自行推断订单状态。
`

// RegisterContextFeatures 注册本讲使用的 Resources 和 Prompts。
func RegisterContextFeatures(server *mcp.Server) {
	registerOrderStatusResource(server)
	registerOrderPrompt(server)
}

// registerOrderStatusResource 注册订单状态说明 Resource。
func registerOrderStatusResource(server *mcp.Server) {
	server.AddResource(
		&mcp.Resource{
			URI:         orderStatusURI,
			Name:        "订单状态说明",
			Description: "解释订单状态的业务含义，不包含实时订单数据。",
			MIMEType:    "text/markdown",
		},
		func(
			ctx context.Context,
			req *mcp.ReadResourceRequest,
		) (*mcp.ReadResourceResult, error) {
			if req.Params.URI != orderStatusURI {
				return nil, mcp.ResourceNotFoundError(req.Params.URI)
			}

			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{
					{
						URI:      orderStatusURI,
						MIMEType: "text/markdown",
						Text:     orderStatusGuide,
					},
				},
			}, nil
		},
	)
}

// registerOrderPrompt 注册订单诊断 Prompt。
func registerOrderPrompt(server *mcp.Server) {
	server.AddPrompt(
		&mcp.Prompt{
			Name:        "diagnose_order",
			Title:       "订单诊断",
			Description: "按照统一流程诊断订单状态和潜在问题。",
			Arguments: []*mcp.PromptArgument{
				{
					Name:        "order_id",
					Title:       "订单 ID",
					Description: "需要诊断的订单 ID，例如 ORD-1001。",
					Required:    true,
				},
			},
		},
		func(
			ctx context.Context,
			req *mcp.GetPromptRequest,
		) (*mcp.GetPromptResult, error) {
			orderID := strings.TrimSpace(req.Params.Arguments["order_id"])
			if orderID == "" {
				return nil, fmt.Errorf("order_id is required")
			}

			// 示例项目中的订单 ID 采用 ORD-数字 的格式。
			// 对 Prompt 参数进行基本校验，避免把任意文本当成订单 ID。
			if !strings.HasPrefix(orderID, "ORD-") {
				return nil, fmt.Errorf("invalid order_id format")
			}

			suffix := strings.TrimPrefix(orderID, "ORD-")
			if suffix == "" {
				return nil, fmt.Errorf("invalid order_id format")
			}
			for _, r := range suffix {
				if r < '0' || r > '9' {
					return nil, fmt.Errorf("invalid order_id format")
				}
			}

			promptText := fmt.Sprintf(`
请诊断订单 %s。

请严格遵循以下工作流程：

1. 使用 get_order Tool 查询订单实时信息。
2. 参考 biz://docs/order-status 中的订单状态说明，
   理解当前订单状态的业务含义。
3. 如果诊断涉及库存问题，可以使用
   inventory_get_item 或 inventory_list_low_stock
   查询相关库存信息。
4. 将结论分成三个部分：
   - 已确认事实：只能引用 Tool 返回的数据。
   - 可能原因：明确标注为推测，不要当成事实。
   - 处理建议：给出可执行的建议，并说明依据。
5. 如果缺少必要信息，应明确指出，不要编造数据。
6. 本次诊断只允许读取数据，不要创建工单、
   取消订单或执行其他写操作。

请使用清晰、简洁的中文回答。
`, orderID)

			return &mcp.GetPromptResult{
				Description: "订单诊断工作流程",
				Messages: []*mcp.PromptMessage{
					{
						Role: "user",
						Content: &mcp.TextContent{
							Text: promptText,
						},
					},
				},
			}, nil
		},
	)
}
