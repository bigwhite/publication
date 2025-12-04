package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bigwhite/my-issue2md/internal/config"
)

const (
	webName    = "issue2mdweb"
	webVersion = "1.0.0"
	defaultPort = "8080"
)

func main() {
	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal, gracefully shutting down...")
		cancel()
	}()

	// 加载配置
	cfg := config.DefaultConfig()
	cfg.LoadFromEnv()

	// 获取端口
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// 启动Web服务器
	if err := runWebServer(ctx, port, cfg); err != nil {
		log.Fatalf("Web server failed: %v", err)
	}
}

// runWebServer 运行Web服务器
func runWebServer(ctx context.Context, port string, cfg *config.Config) error {
	mux := http.NewServeMux()

	// 设置路由
	setupRoutes(mux, cfg)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 启动服务器
	go func() {
		log.Printf("%s v%s starting on port %s", webName, webVersion, port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	// 等待上下文取消
	<-ctx.Done()

	// 优雅关闭
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Println("Web server stopped gracefully")
	return nil
}

// setupRoutes 设置路由
func setupRoutes(mux *http.ServeMux, cfg *config.Config) {
	// 健康检查
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","service":"%s","version":"%s"}`, webName, webVersion)
	})

	// 首页
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>%s v%s</title>
    <meta charset="utf-8">
</head>
<body>
    <h1>%s v%s</h1>
    <p>GitHub Issue to Markdown Converter Web Service</p>
    <p>Phase 1 Foundation completed successfully!</p>
    <ul>
        <li><a href="/health">Health Check</a></li>
        <li><a href="/api/v1/convert">API Documentation</a></li>
    </ul>
</body>
</html>`, webName, webVersion, webName, webVersion)
	})

	// API路由
	mux.HandleFunc("/api/v1/convert", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, `{"error":"API endpoint not yet implemented","status":"Phase 1 in progress"}`)
	})
}