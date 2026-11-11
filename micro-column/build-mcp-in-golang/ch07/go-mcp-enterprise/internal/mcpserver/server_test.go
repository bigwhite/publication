package mcpserver

import (
	"context"
	"encoding/json"
	"log/slog"
	"testing"

	"example.com/go-mcp-enterprise/internal/orders"
	"example.com/go-mcp-enterprise/internal/tickets"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestGetOrderTool(t *testing.T) {
	ctx := context.Background()

	// 创建订单业务服务
	orderService := orders.NewService(
		orders.NewMemoryRepository(),
	)

	// 创建工单业务服务
	ticketService := tickets.NewService(
		orderService,
		tickets.NewMemoryRepository(),
	)

	// 创建 MCP Server
	server := New(
		orderService,
		ticketService,
		slog.Default(),
	)

	// 创建内存传输
	clientTransport, serverTransport := mcp.NewInMemoryTransports()

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

	// 1. 发现 Tools
	tools, err := clientSession.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("ListTools() error = %v", err)
	}

	if len(tools.Tools) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(tools.Tools))
	}

	toolNames := make(map[string]bool)
	for _, tool := range tools.Tools {
		toolNames[tool.Name] = true
	}

	if !toolNames["get_order"] {
		t.Fatal("get_order tool not found")
	}

	if !toolNames["create_ticket"] {
		t.Fatal("create_ticket tool not found")
	}

	// 2. 调用 get_order
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
		t.Fatalf("CallTool(get_order) error = %v", err)
	}

	if result.IsError {
		t.Fatalf("unexpected tool error: %+v", result.Content)
	}

	// 3. 检查 get_order 结构化结果
	data, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatalf("marshal structured content: %v", err)
	}

	var orderOutput GetOrderOutput
	if err := json.Unmarshal(data, &orderOutput); err != nil {
		t.Fatalf("unmarshal order output: %v", err)
	}

	if orderOutput.ID != "ORD-1001" {
		t.Fatalf("unexpected order ID: %s", orderOutput.ID)
	}

	if orderOutput.Status != "paid" {
		t.Fatalf("unexpected order status: %s", orderOutput.Status)
	}

	if orderOutput.TotalAmount != 12990 {
		t.Fatalf("unexpected total amount: %d", orderOutput.TotalAmount)
	}

	// 4. 调用 create_ticket
	ticketResult, err := clientSession.CallTool(
		ctx,
		&mcp.CallToolParams{
			Name: "create_ticket",
			Arguments: map[string]any{
				"order_id":        "ORD-1001",
				"title":           "商品无法正常使用",
				"description":     "收到商品后无法启动。",
				"priority":        "normal",
				"idempotency_key": "test-create-ticket-001",
			},
		},
	)
	if err != nil {
		t.Fatalf("CallTool(create_ticket) error = %v", err)
	}

	if ticketResult.IsError {
		t.Fatalf("unexpected create_ticket error: %+v", ticketResult.Content)
	}

	var ticketOutput CreateTicketOutput
	ticketData, err := json.Marshal(ticketResult.StructuredContent)
	if err != nil {
		t.Fatalf("marshal ticket structured content: %v", err)
	}

	if err := json.Unmarshal(ticketData, &ticketOutput); err != nil {
		t.Fatalf("unmarshal ticket output: %v", err)
	}

	if ticketOutput.TicketID == "" {
		t.Fatal("expected a non-empty ticket ID")
	}

	if ticketOutput.OrderID != "ORD-1001" {
		t.Fatalf("unexpected ticket order ID: %s", ticketOutput.OrderID)
	}

	if ticketOutput.Status != "open" {
		t.Fatalf("unexpected ticket status: %s", ticketOutput.Status)
	}

	// 5. 再次使用相同幂等键和参数调用，应该返回同一张工单
	retryResult, err := clientSession.CallTool(
		ctx,
		&mcp.CallToolParams{
			Name: "create_ticket",
			Arguments: map[string]any{
				"order_id":        "ORD-1001",
				"title":           "商品无法正常使用",
				"description":     "收到商品后无法启动。",
				"priority":        "normal",
				"idempotency_key": "test-create-ticket-001",
			},
		},
	)
	if err != nil {
		t.Fatalf("retry CallTool(create_ticket) error = %v", err)
	}

	if retryResult.IsError {
		t.Fatalf("unexpected retry error: %+v", retryResult.Content)
	}

	var retryOutput CreateTicketOutput
	retryData, err := json.Marshal(retryResult.StructuredContent)
	if err != nil {
		t.Fatalf("marshal retry structured content: %v", err)
	}

	if err := json.Unmarshal(retryData, &retryOutput); err != nil {
		t.Fatalf("unmarshal retry output: %v", err)
	}

	if retryOutput.TicketID != ticketOutput.TicketID {
		t.Fatalf(
			"idempotent retry returned different ticket IDs: first=%s retry=%s",
			ticketOutput.TicketID,
			retryOutput.TicketID,
		)
	}
}
