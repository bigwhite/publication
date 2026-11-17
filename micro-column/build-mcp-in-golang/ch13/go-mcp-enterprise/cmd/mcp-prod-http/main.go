package main

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"example.com/go-mcp-enterprise/internal/inventory"
	"example.com/go-mcp-enterprise/internal/mcpserver"
	"example.com/go-mcp-enterprise/internal/orders"
	"example.com/go-mcp-enterprise/internal/runtimeconfig"
	"example.com/go-mcp-enterprise/internal/tickets"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	if err := run(); err != nil {
		slog.Error("MCP server exited", "error", err)
		os.Exit(1)
	}
}

func run() error {
	logger := slog.New(slog.NewJSONHandler(
		os.Stderr,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	))

	cfg, err := runtimeconfig.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger.Info(
		"starting MCP server",
		"environment", cfg.Environment,
		"address", cfg.Addr,
		"mcp_path", cfg.MCPPath,
	)

	// 初始化业务服务。
	orderService := orders.NewService(
		orders.NewMemoryRepository(),
	)

	ticketService := tickets.NewService(
		orderService,
		tickets.NewMemoryRepository(),
	)

	inventoryService := inventory.NewService(
		inventory.NewMemoryRepository(),
	)

	// 复用现有 MCP Server 和业务 Tool。
	server := mcpserver.New(
		orderService,
		ticketService,
		logger,
	)

	mcpserver.RegisterInventoryTools(
		server,
		inventoryService,
	)

	// Streamable HTTP Handler。
	mcpHandler := mcp.NewStreamableHTTPHandler(
		func(r *http.Request) *mcp.Server {
			return server
		},
		&mcp.StreamableHTTPOptions{
			Stateless: true,
		},
	)

	var ready atomic.Bool
	ready.Store(true)

	mux := http.NewServeMux()

	// 存活检查：只判断进程是否仍在运行。
	mux.HandleFunc("/healthz", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if r.Method != http.MethodGet {
			http.Error(
				w,
				"method not allowed",
				http.StatusMethodNotAllowed,
			)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})

	// 就绪检查：退出流程开始后，不再接收新流量。
	mux.HandleFunc("/readyz", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if r.Method != http.MethodGet {
			http.Error(
				w,
				"method not allowed",
				http.StatusMethodNotAllowed,
			)
			return
		}

		if !ready.Load() {
			http.Error(
				w,
				"not ready",
				http.StatusServiceUnavailable,
			)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready\n"))
	})

	// 仅对 MCP Endpoint 应用 Host、Origin 和认证检查。
	mux.Handle(
		cfg.MCPPath,
		securityMiddleware(
			mcpHandler,
			cfg,
		),
	)

	httpServer := &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
		// 不设置过短的 WriteTimeout。
		// Streamable HTTP 可能需要持续写出响应流。
	}

	// 监听 SIGINT 和 SIGTERM。
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	errCh := make(chan error, 1)

	go func() {
		logger.Info(
			"HTTP server listening",
			"address", cfg.Addr,
		)

		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("HTTP server: %w", err)
		}
		return nil

	case <-ctx.Done():
		logger.Info("shutdown signal received")

		// 先将实例标记为不可接收新流量。
		ready.Store(false)

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			cfg.ShutdownTimeout,
		)
		defer cancel()

		logger.Info(
			"shutting down HTTP server",
			"timeout", cfg.ShutdownTimeout,
		)

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			// Shutdown 超时后，尝试强制关闭连接。
			logger.Error(
				"graceful shutdown failed",
				"error", err,
			)

			if closeErr := httpServer.Close(); closeErr != nil {
				logger.Error(
					"force close failed",
					"error", closeErr,
				)
			}

			return fmt.Errorf(
				"graceful shutdown: %w",
				err,
			)
		}

		logger.Info("HTTP server stopped")
		return nil
	}
}

// securityMiddleware 对 MCP Endpoint 进行基础安全检查。
func securityMiddleware(
	next http.Handler,
	cfg runtimeconfig.Config,
) http.Handler {
	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if !allowedHost(r.Host, cfg.AllowedHosts) {
			http.Error(
				w,
				"forbidden host",
				http.StatusForbidden,
			)
			return
		}

		// 没有 Origin 的非浏览器客户端仍然可以访问，
		// 但如果携带 Origin，就必须通过精确 allowlist 检查。
		origin := r.Header.Get("Origin")
		if origin != "" &&
			!allowedOrigin(origin, cfg.AllowedOrigins) {
			http.Error(
				w,
				"forbidden origin",
				http.StatusForbidden,
			)
			return
		}

		// 开发环境允许不配置 Token，便于使用现有 mcp-host。
		// 生产环境在配置加载阶段已经强制要求 Token。
		if cfg.BearerToken != "" &&
			!validBearerToken(
				r.Header.Get("Authorization"),
				cfg.BearerToken,
			) {
			w.Header().Set(
				"WWW-Authenticate",
				"Bearer",
			)
			http.Error(
				w,
				"unauthorized",
				http.StatusUnauthorized,
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func allowedHost(
	host string,
	allowedHosts []string,
) bool {
	for _, allowed := range allowedHosts {
		if strings.EqualFold(host, allowed) {
			return true
		}
	}
	return false
}

func allowedOrigin(
	origin string,
	allowedOrigins []string,
) bool {
	for _, allowed := range allowedOrigins {
		if origin == allowed {
			return true
		}
	}
	return false
}

func validBearerToken(
	header string,
	expected string,
) bool {
	scheme, token, ok := strings.Cut(header, " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") {
		return false
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return false
	}

	return subtle.ConstantTimeCompare(
		[]byte(token),
		[]byte(expected),
	) == 1
}
