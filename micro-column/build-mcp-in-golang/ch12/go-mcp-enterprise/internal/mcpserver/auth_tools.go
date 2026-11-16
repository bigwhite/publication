package mcpserver

import (
	"context"
	"errors"
	"fmt"

	"example.com/go-mcp-enterprise/internal/orders"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// demoOrderTenant 是教学用的对象归属映射。
// 生产环境应从订单数据或可信的业务授权服务中获取租户归属。
var demoOrderTenant = map[string]string{
	"ORD-1001": "acme",
	"ORD-2001": "globex",
}

func canReadOrder(tenantID, orderID string) bool {
	return demoOrderTenant[orderID] == tenantID
}

func RegisterAuthorizedOrderTools(
	server *mcp.Server,
	orderService *orders.Service,
) {
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "get_order",
			Description: "查询当前用户所属租户内的订单。只能访问当前身份有权查看的订单。",
		},
		func(
			ctx context.Context,
			req *mcp.CallToolRequest,
			input GetOrderInput,
		) (*mcp.CallToolResult, GetOrderOutput, error) {
			if req.Extra == nil || req.Extra.TokenInfo == nil {
				return nil, GetOrderOutput{}, errors.New("unauthenticated")
			}

			tokenInfo := req.Extra.TokenInfo

			if !containsScope(tokenInfo.Scopes, "orders:read") {
				return nil, GetOrderOutput{}, errors.New(
					"forbidden: orders:read scope required",
				)
			}

			tenantID, ok := tokenInfo.Extra["tenant_id"].(string)
			if !ok || tenantID == "" {
				return nil, GetOrderOutput{}, errors.New(
					"missing tenant identity",
				)
			}

			// 先检查对象访问范围，再读取订单。
			// tenantID 来自经过验证的 Token，而不是 Tool 参数。
			if !canReadOrder(tenantID, input.OrderID) {
				return nil, GetOrderOutput{}, fmt.Errorf(
					"forbidden: order %s is outside tenant %s",
					input.OrderID,
					tenantID,
				)
			}

			order, err := orderService.GetOrder(ctx, input.OrderID)
			if err != nil {
				return nil, GetOrderOutput{}, err
			}

			return nil, GetOrderOutput{
				ID:          order.ID,
				CustomerID:  order.CustomerID,
				Status:      order.Status,
				TotalAmount: order.TotalAmount,
				Currency:    order.Currency,
				CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
			}, nil
		},
	)
}

func containsScope(scopes []string, required string) bool {
	for _, scope := range scopes {
		if scope == required {
			return true
		}
	}
	return false
}
