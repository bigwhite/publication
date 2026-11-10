package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"example.com/go-mcp-enterprise/internal/mcpserver"
	"example.com/go-mcp-enterprise/internal/orders"
	"example.com/go-mcp-enterprise/internal/tickets"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	logger := slog.New(
		slog.NewTextHandler(os.Stderr, nil),
	)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	orderRepo := orders.NewMemoryRepository()
	orderService := orders.NewService(orderRepo)

	ticketRepo := tickets.NewMemoryRepository()
	ticketService := tickets.NewService(orderService, ticketRepo)

	// 创建 MCP Server，并注入业务 Service
	server := mcpserver.New(
		orderService,
		ticketService,
		logger,
	)

	// 通过 stdio Transport 启动服务
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		if !errors.Is(err, context.Canceled) {
			logger.Error("MCP server stopped", "error", err)
			os.Exit(1)
		}
	}
}
