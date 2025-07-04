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
	"github.com/your_org/shortlink/internal/idgen"            // 导入 idgen 接口包
	"github.com/your_org/shortlink/internal/idgen/simplehash" // 导入 idgen 具体实现
	"github.com/your_org/shortlink/internal/shortener"
	"github.com/your_org/shortlink/internal/storage"        // 导入 storage 接口包
	"github.com/your_org/shortlink/internal/storage/memory" // 导入 storage 具体实现
)

func main() {
	// 1. 简化配置处理 (硬编码)
	c, _ := config.LoadConfig()

	// 2. 使用标准库 log
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting shortlink service", "version", "shortlink-demo2")

	// 3. 初始化依赖 (依赖接口，注入具体实现)
	var storeInstance storage.Store = memory.New() // memory.Store 实现了 storage.Store 接口
	defer func() {
		if err := storeInstance.Close(); err != nil {
			log.Printf("Error closing store: %v\n", err)
		}
	}()

	var idGeneratorInstance idgen.Generator = simplehash.New() // simplehash.Generator 实现了 idgen.Generator 接口

	// 创建 Service, 注入接口实现
	shortenerService, err := shortener.NewService(shortener.Config{
		Store:           storeInstance,
		Generator:       idGeneratorInstance,
		Logger:          nil, // 传递 nil, Service 内部会使用默认 logger
		MaxGenAttempts:  3,
		MinShortCodeLen: 5,
	})
	if err != nil {
		log.Fatalf("Failed to create shortener service: %v\n", err)
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
