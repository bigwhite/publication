package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"example.com/go-mcp-enterprise/internal/mcpserver"
	"example.com/go-mcp-enterprise/internal/orders"
	"example.com/go-mcp-enterprise/internal/resilience"
	"example.com/go-mcp-enterprise/internal/tickets"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()

	logger := slog.New(
		slog.NewTextHandler(os.Stderr, nil),
	)

	orderService := orders.NewService(
		orders.NewMemoryRepository(),
	)

	ticketService := tickets.NewService(
		orderService,
		tickets.NewMemoryRepository(),
	)

	server := mcpserver.New(
		orderService,
		ticketService,
		logger,
	)

	reader := resilience.NewDelayedOrderReader(
		orderService,
		250000000, // 250ms，仅用于本地演示
	)

	gate := resilience.NewGate(2)
	limiter := resilience.NewLimiter(20, 1000000000) // 20 次/秒

	mcpserver.RegisterResilientOrderTools(
		server,
		reader,
		gate,
		limiter,
	)

	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatalf("MCP Server stopped: %v", err)
	}
}
