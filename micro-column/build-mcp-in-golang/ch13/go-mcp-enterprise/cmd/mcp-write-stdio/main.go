package main

import (
	"context"
	"log/slog"
	"os"

	"example.com/go-mcp-enterprise/internal/cancellations"
	"example.com/go-mcp-enterprise/internal/inventory"
	"example.com/go-mcp-enterprise/internal/mcpserver"
	"example.com/go-mcp-enterprise/internal/orders"
	"example.com/go-mcp-enterprise/internal/tickets"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// stdio 模式下，stdout 只能用于 MCP 协议通信。
	logger := slog.New(
		slog.NewTextHandler(os.Stderr, nil),
	)

	// 初始化订单服务。
	orderService := orders.NewService(
		orders.NewMemoryRepository(),
	)

	// 初始化工单服务。
	ticketService := tickets.NewService(
		orderService,
		tickets.NewMemoryRepository(),
	)

	// 初始化库存服务。
	inventoryService := inventory.NewService(
		inventory.NewMemoryRepository(),
	)

	// 初始化取消申请服务。
	// 它依赖订单读取接口，但不会直接修改订单状态。
	cancellationService := cancellations.NewService(
		orderService,
		cancellations.NewMemoryRepository(),
	)

	// 创建已有订单和工单 Tools。
	server := mcpserver.New(
		orderService,
		ticketService,
		logger,
	)

	// 注册库存 Tools。
	mcpserver.RegisterInventoryTools(
		server,
		inventoryService,
	)

	// 注册取消申请 Tool。
	mcpserver.RegisterCancellationTools(
		server,
		cancellationService,
	)

	logger.Info("starting write-operation MCP server over stdio")

	if err := server.Run(
		context.Background(),
		&mcp.StdioTransport{},
	); err != nil {
		logger.Error("MCP server stopped", "error", err)
		os.Exit(1)
	}
}
