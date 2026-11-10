package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx := context.Background()

	// 1. 创建 MCP Client
	client := mcp.NewClient(
		&mcp.Implementation{
			Name:    "enterprise-orders-host",
			Version: "0.1.0",
		},
		nil,
	)

	// 2. 创建 CommandTransport，由它启动 MCP Server 子进程
	transport := &mcp.CommandTransport{
		Command: exec.Command(
			"go",
			"run",
			"./cmd/mcp-server",
		),
	}

	// 3. 建立 Client 与 Server 之间的连接
	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatalf("connect to MCP server failed: %v", err)
	}
	defer session.Close()

	fmt.Println("Connected to MCP Server.")

	// 4. 发现 Server 提供的 Tools
	tools, err := session.ListTools(ctx, nil)
	if err != nil {
		log.Fatalf("list tools failed: %v", err)
	}

	fmt.Println("\nAvailable tools:")

	for _, tool := range tools.Tools {
		fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
	}

	// 5. 调用 get_order Tool
	fmt.Println("\nCalling get_order...")

	result, err := session.CallTool(
		ctx,
		&mcp.CallToolParams{
			Name: "get_order",
			Arguments: map[string]any{
				"order_id": "ORD-1001",
			},
		},
	)
	if err != nil {
		log.Fatalf("call tool failed: %v", err)
	}

	if result.IsError {
		fmt.Println("Tool execution failed:")
		for _, content := range result.Content {
			if text, ok := content.(*mcp.TextContent); ok {
				fmt.Println(text.Text)
			}
		}
		return
	}

	// 6. 打印结构化结果
	output, err := json.MarshalIndent(
		result.StructuredContent,
		"",
		"  ",
	)
	if err != nil {
		log.Fatalf("marshal tool result failed: %v", err)
	}

	fmt.Println("\nTool result:")
	fmt.Println(string(output))

	// 7. 调用 create_ticket Tool
	result, err = session.CallTool(
		ctx,
		&mcp.CallToolParams{
			Name: "create_ticket",
			Arguments: map[string]any{
				"order_id":        "ORD-1001",
				"title":           "商品无法正常使用",
				"description":     "收到商品后无法启动，已经尝试重新插拔电源。",
				"priority":        "normal",
				"idempotency_key": "host-demo-001",
			},
		},
	)
	if err != nil {
		log.Fatalf("call create_ticket failed: %v", err)
	}

	if result.IsError {
		log.Fatalf("create_ticket returned an error: %+v", result.Content)
	}

	data, err := json.MarshalIndent(
		result.StructuredContent,
		"",
		"  ",
	)
	if err != nil {
		log.Fatalf("marshal result failed: %v", err)
	}

	fmt.Println(string(data))
}
