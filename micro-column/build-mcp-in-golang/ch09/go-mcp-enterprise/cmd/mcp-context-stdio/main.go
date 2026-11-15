package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"example.com/go-mcp-enterprise/internal/inventory"
	"example.com/go-mcp-enterprise/internal/mcpserver"
	"example.com/go-mcp-enterprise/internal/orders"
	"example.com/go-mcp-enterprise/internal/tickets"
)

func main() {
	// stdio 模式下，日志必须写入 stderr，
	// stdout 留给 MCP 协议通信使用。
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

	// 创建包含订单和工单 Tools 的 MCP Server。
	server := mcpserver.New(
		orderService,
		ticketService,
		logger,
	)

	// 注册第 07 讲中的库存 Tools。
	mcpserver.RegisterInventoryTools(
		server,
		inventoryService,
	)

	// 注册本讲的 Resources 和 Prompts。
	mcpserver.RegisterContextFeatures(server)

	logger.Info("starting enterprise MCP server over stdio")

	if err := server.Run(
		context.Background(),
		&mcp.StdioTransport{},
	); err != nil {
		logger.Error("MCP server stopped", "error", err)
		os.Exit(1)
	}
}
