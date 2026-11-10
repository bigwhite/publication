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

	// 1. 创建 Repository
	repo := orders.NewMemoryRepository()

	// 2. 创建业务 Service
	orderService := orders.NewService(repo)

	// 3. 创建 MCP Server，并注入业务 Service
	server := mcpserver.New(orderService, logger)

	// 4. 通过 stdio Transport 启动服务
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		if !errors.Is(err, context.Canceled) {
			logger.Error("MCP server stopped", "error", err)
			os.Exit(1)
		}
	}
}
