package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/your_org/shortlink/internal/api/http/server"
	"github.com/your_org/shortlink/internal/config"
	"github.com/your_org/shortlink/internal/idgen"
	"github.com/your_org/shortlink/internal/shortener"
	"github.com/your_org/shortlink/internal/storage"
)

func main() {
	// 1. 简化配置处理 (硬编码)
	c, _ := config.LoadConfig()

	// 2. 使用标准库 log
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting shortlink service", "version", "shortlink-demo1")

	// 3. 初始化依赖 (具体类型)
	storeImpl := storage.NewStore()
	defer func() {
		if err := storeImpl.Close(); err != nil {
			log.Printf("Error closing store: %v\n", err)
		}
	}()

	idGenImpl := idgen.NewGenerator()

	// 假设 NewService 接受具体类型，并且可能返回错误（如果依赖为nil）
	// shortenerService, err := shortener.NewService(storeImpl, idGenImpl) // 如果 NewService 返回 error
	shortenerService := shortener.NewService(shortener.Config{
		Store:          storeImpl,
		Generator:      idGenImpl,
		MaxGenAttempts: 3}) // 使用Config结构
	if shortenerService == nil {
		log.Fatalln("Failed to create shortener service due to nil dependencies")
	}

	// 4. 创建并启动 HTTP 服务器 (由 api/http/server.go 负责)
	httpServer := server.New(c.Server.Port, shortenerService) // 将 service 注入

	go func() {
		log.Printf("HTTP server starting on :%s\n", c.Server.Port)
		if err := httpServer.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server ListenAndServe error: %v\n", err)
		}
	}()

	// 5. 实现优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("Received signal %s, shutting down server...\n", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // 增加超时时间
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server exiting")
}
