package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type bearerTransport struct {
	token string
	base  http.RoundTripper
}

func (t bearerTransport) RoundTrip(
	req *http.Request,
) (*http.Response, error) {
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}

	clone := req.Clone(req.Context())
	clone.Header = req.Header.Clone()
	clone.Header.Set("Authorization", "Bearer "+t.token)

	return base.RoundTrip(clone)
}

func main() {
	endpoint := flag.String(
		"endpoint",
		"http://127.0.0.1:8080/mcp",
		"MCP Server endpoint",
	)
	token := flag.String(
		"token",
		"",
		"Bearer access token",
	)
	orderID := flag.String(
		"order-id",
		"ORD-1001",
		"Order ID to query",
	)
	flag.Parse()

	if *token == "" {
		log.Fatal("missing -token")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()

	httpClient := &http.Client{
		Transport: bearerTransport{
			token: *token,
		},
	}

	transport := &mcp.StreamableClientTransport{
		Endpoint:   *endpoint,
		HTTPClient: httpClient,
	}

	client := mcp.NewClient(
		&mcp.Implementation{
			Name:    "mcp-write-host",
			Version: "v0.1.0",
		},
		nil,
	)

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatalf("connect MCP Server: %v", err)
	}
	defer session.Close()

	fmt.Println("Connected to authenticated MCP Server.")

	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "get_order",
		Arguments: map[string]any{
			"order_id": *orderID,
		},
	})
	if err != nil {
		log.Fatalf("call get_order: %v", err)
	}

	if result.IsError {
		fmt.Println("Tool call was rejected by the server.")
		for _, content := range result.Content {
			if text, ok := content.(*mcp.TextContent); ok {
				fmt.Println(text.Text)
			}
		}
		return
	}

	if result.StructuredContent != nil {
		data, err := json.MarshalIndent(
			result.StructuredContent,
			"",
			"  ",
		)
		if err != nil {
			log.Fatalf("marshal result: %v", err)
		}
		fmt.Println(string(data))
		return
	}

	for _, content := range result.Content {
		if text, ok := content.(*mcp.TextContent); ok {
			fmt.Println(text.Text)
		}
	}
}
