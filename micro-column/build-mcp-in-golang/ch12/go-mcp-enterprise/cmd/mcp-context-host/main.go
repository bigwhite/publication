package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	orderStatusURI = "biz://docs/order-status"
	testOrderID    = "ORD-1001"
)

func main() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		2*time.Minute,
	)
	defer cancel()

	// 创建 MCP Client。
	client := mcp.NewClient(
		&mcp.Implementation{
			Name:    "mcp-context-host",
			Version: "v0.1.0",
		},
		nil,
	)

	// 启动 mcp-context-stdio 子进程。
	// CommandTransport 负责通过子进程的 stdin/stdout
	// 与 MCP Server 进行 stdio 通信。
	cmd := exec.CommandContext(
		ctx,
		"go",
		"run",
		"./cmd/mcp-context-stdio",
	)

	transport := &mcp.CommandTransport{
		Command: cmd,
	}

	fmt.Println("Connecting to mcp-context-stdio...")

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatalf("connect MCP server: %v", err)
	}
	defer session.Close()

	fmt.Println("Connected successfully.")

	// 1. 验证 Tools。
	if err := checkTools(ctx, session); err != nil {
		log.Fatalf("check tools: %v", err)
	}

	// 2. 调用一个只读 Tool，验证业务调用仍然正常。
	if err := checkGetOrder(ctx, session); err != nil {
		log.Fatalf("call get_order: %v", err)
	}

	// 3. 验证 Resource 的发现和读取。
	if err := checkOrderStatusResource(ctx, session); err != nil {
		log.Fatalf("check resource: %v", err)
	}

	// 4. 验证 Prompt 的发现和获取。
	if err := checkOrderPrompt(ctx, session); err != nil {
		log.Fatalf("check prompt: %v", err)
	}

	fmt.Println()
	fmt.Println("All MCP context checks passed.")
}

// checkTools 列出 Server 提供的 Tools。
func checkTools(
	ctx context.Context,
	session *mcp.ClientSession,
) error {
	fmt.Println()
	fmt.Println("========== 1. Tools ==========")

	result, err := session.ListTools(ctx, nil)
	if err != nil {
		return fmt.Errorf("list tools: %w", err)
	}

	for _, tool := range result.Tools {
		fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
	}

	fmt.Printf("Total tools: %d\n", len(result.Tools))
	return nil
}

// checkGetOrder 调用已有的 get_order Tool。
func checkGetOrder(
	ctx context.Context,
	session *mcp.ClientSession,
) error {
	fmt.Println()
	fmt.Println("========== 2. Call Tool: get_order ==========")

	result, err := session.CallTool(
		ctx,
		&mcp.CallToolParams{
			Name: "get_order",
			Arguments: map[string]any{
				"order_id": testOrderID,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("call get_order: %w", err)
	}

	if result.IsError {
		return fmt.Errorf("get_order returned an MCP tool error")
	}

	if result.StructuredContent != nil {
		data, err := json.MarshalIndent(
			result.StructuredContent,
			"",
			"  ",
		)
		if err != nil {
			return fmt.Errorf("marshal tool result: %w", err)
		}

		fmt.Println(string(data))
	} else {
		for _, content := range result.Content {
			if text, ok := content.(*mcp.TextContent); ok {
				fmt.Println(text.Text)
			}
		}
	}

	return nil
}

// checkOrderStatusResource 验证订单状态 Resource。
func checkOrderStatusResource(
	ctx context.Context,
	session *mcp.ClientSession,
) error {
	fmt.Println()
	fmt.Println("========== 3. Resource ==========")

	found := false

	// 通过 resources/list 发现 Resource。
	for resource, err := range session.Resources(ctx, nil) {
		if err != nil {
			return fmt.Errorf("list resources: %w", err)
		}

		fmt.Printf(
			"- %s (%s)\n",
			resource.Name,
			resource.URI,
		)

		if resource.URI == orderStatusURI {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("resource %q not found", orderStatusURI)
	}

	// 通过 resources/read 读取 Resource。
	result, err := session.ReadResource(
		ctx,
		&mcp.ReadResourceParams{
			URI: orderStatusURI,
		},
	)
	if err != nil {
		return fmt.Errorf("read resource: %w", err)
	}

	fmt.Printf("Resource URI: %s\n", orderStatusURI)

	for _, content := range result.Contents {
		fmt.Printf("MIME Type: %s\n", content.MIMEType)
		fmt.Println("Resource Content:")
		fmt.Println(content.Text)
	}

	return nil
}

// checkOrderPrompt 验证订单诊断 Prompt。
func checkOrderPrompt(
	ctx context.Context,
	session *mcp.ClientSession,
) error {
	fmt.Println()
	fmt.Println("========== 4. Prompt ==========")

	found := false

	// 通过 prompts/list 发现 Prompt。
	for prompt, err := range session.Prompts(ctx, nil) {
		if err != nil {
			return fmt.Errorf("list prompts: %w", err)
		}

		fmt.Printf("- %s: %s\n", prompt.Name, prompt.Description)

		if prompt.Name == "diagnose_order" {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("prompt diagnose_order not found")
	}

	// 通过 prompts/get 获取填充参数后的 Prompt。
	result, err := session.GetPrompt(
		ctx,
		&mcp.GetPromptParams{
			Name: "diagnose_order",
			Arguments: map[string]string{
				"order_id": testOrderID,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("get prompt: %w", err)
	}

	fmt.Println("Prompt: diagnose_order")

	for _, message := range result.Messages {
		fmt.Printf("Role: %s\n", message.Role)

		switch content := message.Content.(type) {
		case *mcp.TextContent:
			fmt.Println(content.Text)

		case *mcp.EmbeddedResource:
			fmt.Printf(
				"Embedded Resource: %s\n",
				content.Resource.URI,
			)

		default:
			fmt.Printf(
				"Unsupported prompt content type: %T\n",
				content,
			)
		}
	}

	return nil
}
