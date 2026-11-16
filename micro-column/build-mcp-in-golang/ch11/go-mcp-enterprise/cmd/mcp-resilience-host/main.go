package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const testOrderID = "ORD-1001"

func main() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()

	client := mcp.NewClient(
		&mcp.Implementation{
			Name:    "mcp-resilience-host",
			Version: "v0.1.0",
		},
		nil,
	)

	cmd := exec.CommandContext(
		ctx,
		"go",
		"run",
		"./cmd/mcp-resilience-stdio",
	)

	transport := &mcp.CommandTransport{
		Command: cmd,
	}

	fmt.Println("Connecting to resilience MCP Server...")

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		log.Fatalf("connect MCP Server: %v", err)
	}
	defer session.Close()

	fmt.Println("Connected successfully.")

	if err := checkToolDiscovery(ctx, session); err != nil {
		log.Fatalf("discover tool: %v", err)
	}

	if err := checkCancellation(ctx, session); err != nil {
		log.Fatalf("cancellation test: %v", err)
	}

	if err := checkConcurrentCalls(ctx, session); err != nil {
		log.Fatalf("concurrency test: %v", err)
	}

	fmt.Println("\nAll resilience checks completed.")
}

func checkToolDiscovery(
	ctx context.Context,
	session *mcp.ClientSession,
) error {
	result, err := session.ListTools(ctx, nil)
	if err != nil {
		return err
	}

	found := false
	for _, tool := range result.Tools {
		fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
		if tool.Name == "get_order_resilient" {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("get_order_resilient not found")
	}

	return nil
}

func callGetOrder(
	ctx context.Context,
	session *mcp.ClientSession,
) error {
	result, err := session.CallTool(
		ctx,
		&mcp.CallToolParams{
			Name: "get_order_resilient",
			Arguments: map[string]any{
				"order_id": testOrderID,
			},
		},
	)
	if err != nil {
		return err
	}

	if result.IsError {
		return fmt.Errorf("tool returned an error")
	}

	fmt.Println("Order query succeeded.")
	return nil
}

func checkCancellation(
	ctx context.Context,
	session *mcp.ClientSession,
) error {
	fmt.Println("\n========== Cancellation test ==========")

	callCtx, cancel := context.WithTimeout(
		ctx,
		40*time.Millisecond,
	)
	defer cancel()

	err := callGetOrder(callCtx, session)
	if err == nil {
		return fmt.Errorf("expected cancellation or deadline error")
	}

	fmt.Printf("Call ended as expected: %v\n", err)

	// 取消一次调用后，同一个 MCP Session 应仍然能够继续使用。
	fmt.Println("Checking session after cancellation...")

	if err := callGetOrder(ctx, session); err != nil {
		return fmt.Errorf("session unusable after cancellation: %w", err)
	}

	return nil
}

func checkConcurrentCalls(
	ctx context.Context,
	session *mcp.ClientSession,
) error {
	fmt.Println("\n========== Concurrent calls ==========")

	const count = 4

	var wg sync.WaitGroup
	errs := make(chan error, count)

	for i := 0; i < count; i++ {
		wg.Add(1)

		go func(index int) {
			defer wg.Done()

			fmt.Printf("Starting concurrent call %d\n", index+1)

			if err := callGetOrder(ctx, session); err != nil {
				errs <- fmt.Errorf("call %d: %w", index+1, err)
			}
		}(i)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		return err
	}

	fmt.Printf("All %d concurrent calls succeeded.\n", count)
	return nil
}
