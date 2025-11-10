package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"context"
	"time"
	"fmt"
	"encoding/json"
	"strings"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	
	"github.com/bigwhite/shortlink/pkg/handler"
	"github.com/bigwhite/shortlink/pkg/repository/postgres"
	redis_repo "github.com/bigwhite/shortlink/pkg/repository/redis"
	"github.com/bigwhite/shortlink/pkg/service"
)

func main() {
	// --- 配置加载 ---
	// 在真实应用中，会使用 viper 等库从文件或环境变量加载
	pgDSN := os.Getenv("DB_DSN")
	if pgDSN == "" {
		log.Fatal("环境变量 DB_DSN 未设置")
	}
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("环境变量 REDIS_ADDR 未设置")
	}

	// --- 依赖初始化 ---
	var dbPool *sql.DB
	var err error

	const maxDBRetries = 10
	const dbRetryDelay = 3 * time.Second

	log.Println("正在连接到 PostgreSQL...")
	for i := 0; i < maxDBRetries; i++ {
		dbPool, err = sql.Open("postgres", pgDSN)
		if err == nil {
			if err = dbPool.Ping(); err == nil {
				log.Println("成功连接到 PostgreSQL！")
				break // 连接成功，跳出循环
			}
		}
		log.Printf("连接 PG 失败 (尝试 #%d/%d): %v", i+1, maxDBRetries, err)
		if i < maxDBRetries-1 {
			time.Sleep(dbRetryDelay)
		}
	}
	if err != nil {
		log.Fatalf("在 %d 次尝试后，仍无法连接到 PostgreSQL: %s", maxDBRetries, err)
	}
	// 在 main 函数结束时关闭连接池
	defer dbPool.Close()

	redisClient := redis.NewClient(&redis.Options{Addr: redisAddr})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("无法连接 Redis: %s", err)
	}
	
	// --- 组装应用 ---
	linkRepo := postgres.NewPgLinkRepository(dbPool)
	linkCache := redis_repo.NewRedisLinkCache(redisClient)
	shortenerSvc := service.NewShortenerService(linkRepo, linkCache)
	linkHandler := handler.NewLinkHandler(shortenerSvc)

	// --- 路由与服务器启动 ---
	mux := http.NewServeMux()
	mux.HandleFunc("/api/links", linkHandler.CreateLink)

	// (为其他 handler 添加路由...)
	// 为统计接口注册一个更具体的路径
	mux.HandleFunc("/api/links/", func(w http.ResponseWriter, r *http.Request) {
		// 这是一个简单的、基于路径前缀的“微路由”
		if strings.HasSuffix(r.URL.Path, "/stats") {
			linkHandler.GetStats(w, r)
			return
		}
		http.NotFound(w, r)
	})

	// 将重定向处理器注册到根路径，处理所有其他请求
	mux.HandleFunc("/", linkHandler.Redirect)


	// --- 新增 /healthz 端点 ---
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()

		// 检查数据库连接
		if err := dbPool.PingContext(ctx); err != nil {
			http.Error(w, fmt.Sprintf("db ping failed: %v", err), http.StatusServiceUnavailable)
			return
		}

		// 检查 Redis 连接
		if err := redisClient.Ping(ctx).Err(); err != nil {
			http.Error(w, fmt.Sprintf("redis ping failed: %v", err), http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	log.Println("Shortlink Service 启动于 :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("服务器启动失败: %s", err)
	}
}
