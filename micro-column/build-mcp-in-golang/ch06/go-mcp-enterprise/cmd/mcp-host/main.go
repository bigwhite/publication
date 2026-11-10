package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	transportType := flag.String(
		"transport",
		"stdio",
		"MCP transport: stdio or http",
	)

	endpoint := flag.String(
		"endpoint",
		"http://127.0.0.1:8080/mcp",
		"Streamable HTTP MCP endpoint",
	)

	flag.Parse()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		60*time.Second,
	)
	defer cancel()

	client := mcp.NewClient(
		&mcp.Implementation{
			Name:    "enterprise-orders-host",
			Version: "0.1.0",
		},
		nil,
	)

	var transport mcp.Transport

	switch *transportType {
	case "stdio":
		fmt.Println("Connecting through stdio...")

		cmd := exec.Command(
			"go",
			"run",
			"./cmd/mcp-stdio",
		)

		// 子进程 stderr 用于输出日志和诊断信息。
		// stdout 必须保留给 MCP 协议通信。
		cmd.Stderr = os.Stderr

		transport = &mcp.CommandTransport{
			Command: cmd,
		}

	case "http":
		fmt.Printf(
			"Connecting through Streamable HTTP: %s\n",
			*endpoint,
		)

		transport = &mcp.StreamableClientTransport{
			Endpoint: *endpoint,
		}

	default:
		log.Fatalf(
			"unsupported transport %q: use stdio or http",
			*transportType,
		)
	}

	session, err := client.Connect(
		ctx,
		transport,
		nil,
	)
	if err != nil {
		log.Fatalf("connect MCP server: %v", err)
	}

	defer func() {
		if err := session.Close(); err != nil {
			fmt.Fprintf(
				os.Stderr,
				"close MCP session: %v\n",
				err,
			)
		}
	}()

	fmt.Println("Connected to MCP server.")

	// 发现服务端提供的 Tools。
	tools, err := session.ListTools(ctx, nil)
	if err != nil {
		log.Fatalf("list tools: %v", err)
	}

	fmt.Println("\nAvailable tools:")
	for _, tool := range tools.Tools {
		fmt.Printf(
			"- %s: %s\n",
			tool.Name,
			tool.Description,
		)
	}

	// 调用 get_order。
	fmt.Println("\nCalling get_order...")

	orderResult, err := session.CallTool(
		ctx,
		&mcp.CallToolParams{
			Name: "get_order",
			Arguments: map[string]any{
				"order_id": "ORD-1001",
			},
		},
	)
	if err != nil {
		log.Fatalf("call get_order: %v", err)
	}

	if orderResult.IsError {
		log.Fatalf("get_order returned an MCP tool error")
	}

	printResult("get_order", orderResult)

	// 调用 create_ticket。
	fmt.Println("\nCalling create_ticket...")

	ticketResult, err := session.CallTool(
		ctx,
		&mcp.CallToolParams{
			Name: "create_ticket",
			Arguments: map[string]any{
				"order_id":        "ORD-1001",
				"title":           "Order delivery verification",
				"description":     "Please verify the delivery status of this order.",
				"priority":        "normal",
				"idempotency_key": "lesson06-host-demo-001",
			},
		},
	)
	if err != nil {
		log.Fatalf("call create_ticket: %v", err)
	}

	if ticketResult.IsError {
		log.Fatalf("create_ticket returned an MCP tool error")
	}

	printResult("create_ticket", ticketResult)

	fmt.Println("\nMCP interaction completed successfully.")
}

func printResult(
	name string,
	result *mcp.CallToolResult,
) {
	data, err := json.MarshalIndent(
		result.StructuredContent,
		"",
		"  ",
	)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"marshal %s result: %v\n",
			name,
			err,
		)
		return
	}

	fmt.Printf("%s result:\n%s\n", name, data)
}
