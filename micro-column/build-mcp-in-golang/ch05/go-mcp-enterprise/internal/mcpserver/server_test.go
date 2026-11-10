package mcpserver

import (
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"example.com/go-mcp-enterprise/internal/orders"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestGetOrderTool(t *testing.T) {
	ctx := context.Background()

	// 创建业务服务
	service := orders.NewService(
		orders.NewMemoryRepository(),
	)

	// 创建 MCP Server
	server := New(service, slog.Default())

	// 创建内存传输
	clientTransport, serverTransport :=
		mcp.NewInMemoryTransports()

	// 连接 Server
	serverSession, err := server.Connect(
		ctx,
		serverTransport,
		nil,
	)
	if err != nil {
		t.Fatalf("server.Connect() error = %v", err)
	}

	// 创建并连接 Client
	client := mcp.NewClient(
		&mcp.Implementation{
			Name:    "test-client",
			Version: "0.1.0",
		},
		nil,
	)

	clientSession, err := client.Connect(
		ctx,
		clientTransport,
		nil,
	)
	if err != nil {
		t.Fatalf("client.Connect() error = %v", err)
	}

	t.Cleanup(func() {
		_ = clientSession.Close()
		_ = serverSession.Wait()
	})

	// 1. 发现 Tool
	tools, err := clientSession.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("ListTools() error = %v", err)
	}

	if len(tools.Tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(tools.Tools))
	}

	if tools.Tools[0].Name != "get_order" {
		t.Fatalf("unexpected tool name: %s", tools.Tools[0].Name)
	}

	// 2. 调用 Tool
	result, err := clientSession.CallTool(
		ctx,
		&mcp.CallToolParams{
			Name: "get_order",
			Arguments: map[string]any{
				"order_id": "ORD-1001",
			},
		},
	)
	if err != nil {
		t.Fatalf("CallTool() error = %v", err)
	}

	if result.IsError {
		t.Fatalf("unexpected tool error: %+v", result.Content)
	}

	// 3. 检查结构化结果
	data, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatalf("marshal structured content: %v", err)
	}

	var output GetOrderOutput
	if err := json.Unmarshal(data, &output); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}

	if output.ID != "ORD-1001" {
		t.Fatalf("unexpected order ID: %s", output.ID)
	}

	if output.Status != "paid" {
		t.Fatalf("unexpected order status: %s", output.Status)
	}

	if output.TotalAmount != 12990 {
		t.Fatalf("unexpected total amount: %d", output.TotalAmount)
	}

	// 4. 测试订单不存在
	missing, err := clientSession.CallTool(
		ctx,
		&mcp.CallToolParams{
			Name: "get_order",
			Arguments: map[string]any{
				"order_id": "ORD-404",
			},
		},
	)
	if err != nil {
		t.Fatalf("CallTool() error = %v", err)
	}

	if !missing.IsError {
		t.Fatal("expected a tool error for missing order")
	}
}
