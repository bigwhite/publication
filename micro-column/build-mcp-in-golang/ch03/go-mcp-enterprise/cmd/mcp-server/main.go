package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"example.com/go-mcp-enterprise/internal/mcpserver"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// 所有日志输出到 stderr，不能污染 stdio 协议通道。
	logger := slog.New(
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	// 接收操作系统退出信号，触发优雅关闭。
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	server := mcpserver.New(logger)

	logger.Info("starting MCP server",
		"name", "go-mcp-enterprise",
		"version", "0.1.0",
	)

	// 当前先使用 stdio Transport。
	err := server.Run(ctx, &mcp.StdioTransport{})
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Printf("MCP server stopped with error: %v", err)
		os.Exit(1)
	}

	logger.Info("MCP server stopped")
}
