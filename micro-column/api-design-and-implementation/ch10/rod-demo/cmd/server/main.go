package main

import (
	"rod-demo/internal/middleware"
	"rod-demo/internal/router"
	"rod-demo/pkg/limiter"
	"rod-demo/pkg/redis"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化 Redis
	redis.Init()
	rateLimiter := limiter.NewLimiter(redis.Client)

	r := gin.Default()

	// 定义限流规则
	// 例如：每分钟 10 次，允许突发 10 次
	generalLimit := limiter.LimitDefinition{
		Rate:   10,
		Period: time.Minute,
		Burst:  10,
	}

	// 注册全局错误处理中间件
	r.Use(middleware.ErrorHandler())

	// API 版本控制（这也是一种资源层级）
	v1 := r.Group("/api/v1")

	// 使用幂等性中间件
	// 通常只针对非 GET/DELETE 请求启用
	v1.Use(middleware.Idempotency())

	v1.Use(middleware.RateLimit(rateLimiter, middleware.IPKeyStrategy, generalLimit))

	router.SetupRoutes(r)

	r.Run(":8080")
}
