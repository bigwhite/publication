package main

import (
	"agentic-discovery-demo/registry"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	// 在系统启动时，硬编码注册几个能力作为演示
	// 实际工程中，这些信息可以通过解析代码注释或专门的配置文件自动生成

	registry.Global.Register(registry.Capability{
		Name:          "fetch_order",
		Category:      registry.CategoryAcquire,
		Description:   "Retrieves details of a specific order. Idempotent and safe to call.",
		Endpoint:      "/agentic/v1/fetch/order",
		Preconditions: []string{"Requires a valid order_id"},
		RequiredScope: "read:orders",
	})

	registry.Global.Register(registry.Capability{
		Name:        "refund_order",
		Category:    registry.CategoryTransact,
		Description: "Initiates a refund for a completed order. HIGH RISK ACTION.",
		Endpoint:    "/agentic/v1/refund/order",
		Preconditions: []string{
			"Order status MUST be 'PAID' or 'SHIPPED'",
			"Refund amount MUST NOT exceed original order amount",
		},
		RequiredScope: "write:finance:refunds",
	})

	registry.Global.Register(registry.Capability{
		Name:          "summarize_document",
		Category:      registry.CategoryCompute,
		Description:   "Analyzes a document and returns a concise summary.",
		Endpoint:      "/agentic/v1/summarize/document",
		Preconditions: []string{"Document must exist in the knowledge base"},
		RequiredScope: "read:documents",
	})
}

func main() {
	r := gin.Default()

	agenticAPI := r.Group("/agentic/v1")
	{
		// 核心：实现动态能力发现接口
		agenticAPI.POST("/discover/actions", DiscoverActions)

		// 下面是具体的执行端点（略，参见上一讲代码）
		// agenticAPI.POST("/fetch/order", ...)
		// agenticAPI.POST("/refund/order", ...)
	}

	log.Println("Agentic API Discovery Server running on http://localhost:8080")
	r.Run(":8080")
}

// DiscoverActions 是 AI 探索系统能力的入口
func DiscoverActions(c *gin.Context) {
	// 获取系统中所有已注册的能力
	allCaps := registry.Global.ListAll()

	// 我们可以根据请求中的参数进行过滤，比如 AI 只想查看 "TRANSACT" 类的能力
	// 这里为了演示简便，返回全部

	c.JSON(http.StatusOK, gin.H{
		"action":       "DISCOVER",
		"status":       "SUCCESS",
		"message":      "Capabilities discovered successfully.",
		"total_count":  len(allCaps),
		"capabilities": allCaps,
	})
}
