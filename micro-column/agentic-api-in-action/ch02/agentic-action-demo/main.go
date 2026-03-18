package main

import (
	"agentic-action/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建一个默认的 Gin 引擎
	r := gin.Default()

	// 我们可以定义一个统一的中间件，专门拦截并记录 Agentic API 的调用意图
	r.Use(IntentLoggerMiddleware())

	// Agentic API 路由组
	agenticAPI := r.Group("/agentic/v1")
	{
		// 1. Acquire (获取) 分类
		// 语义：安全提取数据。AI 知道这个操作是幂等的。
		// 传统写法: r.GET("/orders/:id", handlers.GetOrder)
		// Agentic 写法:
		agenticAPI.POST("/fetch/order", handlers.FetchOrder)
		agenticAPI.POST("/search/orders", handlers.SearchOrders)

		// 2. Transact (交易) 分类
		// 语义：高危状态变更。AI 会警惕此类操作。
		// 传统写法: r.POST("/refunds", handlers.CreateRefund)
		// Agentic 写法:
		agenticAPI.POST("/refund/order", handlers.RefundOrder)
	}

	log.Println("Agentic API Server running on http://localhost:8080")
	r.Run(":8080")
}

// IntentLoggerMiddleware 记录 AI 的调用意图，方便后续审计 (DBOM)
func IntentLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求路径，这就是 AI 的明确意图
		log.Printf("[AGENT INTENT] AI is attempting to execute: %s", c.Request.URL.Path)
		c.Next()
	}
}
