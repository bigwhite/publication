package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		2*time.Minute,
	)
	defer cancel()

	client := mcp.NewClient(
		&mcp.Implementation{
			Name:    "mcp-observe-host",
			Version: "1.0.0",
		},
		nil,
	)

	// CommandTransport 自动启动 MCP Server 子进程。
	cmd := exec.CommandContext(
		ctx,
		"go",
		"run",
		"./cmd/mcp-observe-stdio",
	)
	cmd.Stderr = os.Stderr

	transport := &mcp.CommandTransport{
		Command: cmd,
	}

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		return fmt.Errorf("connect MCP server: %w", err)
	}
	defer session.Close()

	// 1. 发现 Tool。
	list, err := session.ListTools(ctx, nil)
	if err != nil {
		return fmt.Errorf("list tools: %w", err)
	}

	found := false
	for _, tool := range list.Tools {
		fmt.Println("discovered tool:", tool.Name)
		if tool.Name == "get_order_observed" {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("get_order_observed was not registered")
	}

	// 2. 查询存在的订单。
	successResult, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "get_order_observed",
		Arguments: map[string]any{
			"order_id": "ORD-1001",
		},
	})
	if err != nil {
		return fmt.Errorf("call existing order: %w", err)
	}
	if successResult.IsError {
		return fmt.Errorf("existing order unexpectedly returned an error")
	}

	fmt.Println("ORD-1001 lookup: success")

	// 3. 查询不存在的订单，制造一次可预期的 Tool 错误。
	errorResult, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "get_order_observed",
		Arguments: map[string]any{
			"order_id": "ORD-404",
		},
	})
	if err != nil {
		return fmt.Errorf("call missing order: %w", err)
	}
	if !errorResult.IsError {
		return fmt.Errorf("missing order should return a Tool error")
	}

	fmt.Println("ORD-404 lookup: expected Tool error")

	// 4. 读取 MCP Server 暴露的指标。
	metricsURL := "http://127.0.0.1:9091/metrics"

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := httpClient.Get(metricsURL)
	if err != nil {
		return fmt.Errorf("read metrics endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("metrics endpoint returned HTTP %d", resp.StatusCode)
	}

	fmt.Println("\n--- MCP Server Metrics ---")

	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		return fmt.Errorf("read metrics response: %w", err)
	}

	return nil
}
