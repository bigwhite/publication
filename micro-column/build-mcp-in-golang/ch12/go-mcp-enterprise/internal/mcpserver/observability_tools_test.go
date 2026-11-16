package mcpserver

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"example.com/go-mcp-enterprise/internal/observability"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestObservedOrderTool(t *testing.T) {
	ctx := context.Background()

	logger := slog.New(
		slog.NewTextHandler(io.Discard, nil),
	)
	metrics := observability.NewMetrics()

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "test-server",
			Version: "1.0.0",
		},
		nil,
	)

	lookup := func(_ context.Context, orderID string) error {
		if orderID == "ORD-1001" {
			return nil
		}
		return errors.New("order not found")
	}

	RegisterObservedOrderTool(server, lookup, logger, metrics)

	clientTransport, serverTransport := mcp.NewInMemoryTransports()

	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer serverSession.Close()

	client := mcp.NewClient(
		&mcp.Implementation{
			Name:    "test-client",
			Version: "1.0.0",
		},
		nil,
	)

	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer clientSession.Close()

	result, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "get_order_observed",
		Arguments: map[string]any{
			"order_id": "ORD-1001",
		},
	})
	if err != nil {
		t.Fatalf("CallTool() error = %v", err)
	}
	if result.IsError {
		t.Fatalf("expected success, got error result: %+v", result.Content)
	}

	result, err = clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "get_order_observed",
		Arguments: map[string]any{
			"order_id": "ORD-404",
		},
	})
	if err != nil {
		t.Fatalf("CallTool() error = %v", err)
	}
	if !result.IsError {
		t.Fatal("expected an error result for missing order")
	}

	snapshot := metrics.Snapshot()["get_order_observed"]

	if snapshot.Calls != 2 {
		t.Fatalf("Calls = %d, want 2", snapshot.Calls)
	}
	if snapshot.Errors != 1 {
		t.Fatalf("Errors = %d, want 1", snapshot.Errors)
	}
}
