package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"example.com/go-mcp-enterprise/internal/inventory"
	"example.com/go-mcp-enterprise/internal/mcpserver"
	"example.com/go-mcp-enterprise/internal/orders"
	"example.com/go-mcp-enterprise/internal/tickets"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
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

	inventoryService := inventory.NewService(
		inventory.NewMemoryRepository(),
	)

	server := mcpserver.New(
		orderService,
		ticketService,
		logger,
	)

	mcpserver.RegisterInventoryTools(
		server,
		inventoryService,
	)

	if err := server.Run(
		context.Background(),
		&mcp.StdioTransport{},
	); err != nil {
		log.Fatal(err)
	}
}
