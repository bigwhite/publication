package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

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

	server := mcpserver.New(
		orderService,
		ticketService,
		logger,
	)

	handler := mcp.NewStreamableHTTPHandler(
		func(r *http.Request) *mcp.Server {
			return server
		},
		&mcp.StreamableHTTPOptions{
			Stateless: true,
		},
	)

	mux := http.NewServeMux()
	mux.Handle("/mcp", handler)

	httpServer := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Info(
		"MCP HTTP server listening",
		"addr", httpServer.Addr,
		"path", "/mcp",
	)

	if err := httpServer.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
