package mcpserver

import (
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestContextFeatures(t *testing.T) {
	ctx := context.Background()

	// 创建只包含本讲 Resources 和 Prompts 的测试 Server。
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "enterprise-test-server",
			Version: "v0.1.0",
		},
		nil,
	)

	RegisterContextFeatures(server)

	client := mcp.NewClient(
		&mcp.Implementation{
			Name:    "enterprise-test-client",
			Version: "v0.1.0",
		},
		nil,
	)

	// 使用 SDK 的内存传输建立真实 MCP 通信。
	serverTransport, clientTransport := mcp.NewInMemoryTransports()

	if _, err := server.Connect(ctx, serverTransport, nil); err != nil {
		t.Fatalf("connect server: %v", err)
	}

	session, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("connect client: %v", err)
	}
	defer session.Close()

	t.Run("list and read order status resource", func(t *testing.T) {
		found := false

		for resource, err := range session.Resources(ctx, nil) {
			if err != nil {
				t.Fatalf("list resources: %v", err)
			}

			if resource.URI == orderStatusURI {
				found = true
				break
			}
		}

		if !found {
			t.Fatalf("resource %q not found", orderStatusURI)
		}

		result, err := session.ReadResource(
			ctx,
			&mcp.ReadResourceParams{
				URI: orderStatusURI,
			},
		)
		if err != nil {
			t.Fatalf("read resource: %v", err)
		}

		if len(result.Contents) != 1 {
			t.Fatalf("expected 1 resource content, got %d", len(result.Contents))
		}

		content := result.Contents[0]
		if content.URI != orderStatusURI {
			t.Fatalf("unexpected resource URI: %s", content.URI)
		}

		if !strings.Contains(content.Text, "paid") {
			t.Fatalf("resource content does not contain order status")
		}
	})

	t.Run("list and get order prompt", func(t *testing.T) {
		found := false

		for prompt, err := range session.Prompts(ctx, nil) {
			if err != nil {
				t.Fatalf("list prompts: %v", err)
			}

			if prompt.Name == "diagnose_order" {
				found = true
				break
			}
		}

		if !found {
			t.Fatal("prompt diagnose_order not found")
		}

		result, err := session.GetPrompt(
			ctx,
			&mcp.GetPromptParams{
				Name: "diagnose_order",
				Arguments: map[string]string{
					"order_id": "ORD-1001",
				},
			},
		)
		if err != nil {
			t.Fatalf("get prompt: %v", err)
		}

		if len(result.Messages) != 1 {
			t.Fatalf("expected 1 prompt message, got %d", len(result.Messages))
		}

		content, ok := result.Messages[0].Content.(*mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", result.Messages[0].Content)
		}

		if !strings.Contains(content.Text, "ORD-1001") {
			t.Fatal("prompt does not contain the supplied order ID")
		}

		if !strings.Contains(content.Text, "get_order") {
			t.Fatal("prompt does not mention get_order")
		}
	})
}
