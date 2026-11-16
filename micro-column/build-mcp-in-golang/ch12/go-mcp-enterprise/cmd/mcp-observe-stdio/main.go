package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/go-mcp-enterprise/internal/mcpserver"
	"example.com/go-mcp-enterprise/internal/observability"
	"example.com/go-mcp-enterprise/internal/orders"
	"example.com/go-mcp-enterprise/internal/tickets"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// stdio Server 的普通日志必须写入 stderr。
	logger := slog.New(
		slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	metrics := observability.NewMetrics()

	metricsAddr := os.Getenv("MCP_METRICS_ADDR")
	if metricsAddr == "" {
		metricsAddr = "127.0.0.1:9091"
	}

	// 先监听指标端口，确保 Host 调用 Tool 后可以读取指标。
	listener, err := net.Listen("tcp", metricsAddr)
	if err != nil {
		return fmt.Errorf("listen metrics endpoint: %w", err)
	}

	metricsServer := &http.Server{
		Handler:           metrics.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if err := metricsServer.Serve(listener); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			logger.Error("metrics server stopped", "error", err)
		}
	}()

	defer func() {
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		if err := metricsServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("shutdown metrics server", "error", err)
		}
	}()

	// 初始化已有的订单和工单业务服务。
	orderService := orders.NewService(
		orders.NewMemoryRepository(),
	)

	ticketService := tickets.NewService(
		orderService,
		tickets.NewMemoryRepository(),
	)

	// 复用现有 MCP Server，不修改 mcpserver.New 的签名。
	server := mcpserver.New(
		orderService,
		ticketService,
		logger,
	)

	// 通过已有订单 Service 实际查询订单。
	lookup := func(ctx context.Context, orderID string) error {
		_, err := orderService.GetOrder(ctx, orderID)
		return err
	}

	mcpserver.RegisterObservedOrderTool(
		server,
		lookup,
		logger,
		metrics,
	)

	logger.Info(
		"observability demo server starting",
		"metrics_addr",
		metricsAddr,
	)

	// stdout 只用于 MCP stdio 协议通信。
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		return fmt.Errorf("run MCP server: %w", err)
	}

	return nil
}
