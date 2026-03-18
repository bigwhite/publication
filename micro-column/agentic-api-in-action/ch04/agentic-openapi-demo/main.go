package main

import (
	"agentic-openapi-demo/spec"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	// 在应用启动时，硬编码注册一个典型的 Agentic 接口到 OpenAPI 规范中
	// 实际工程中，你可以编写中间件或反射机制，在绑定路由的同时自动提取这些信息

	spec.GlobalSpec.RegisterRoute("/agentic/v1/users/{id}/status", "put", spec.Operation{
		Summary:     "Update user status",
		OperationID: "updateUserStatus",
		// --- 注入 Agentic 元数据 ---
		AgenticMeta: spec.AgenticMeta{
			Action: "TRANSACT",
			Preconditions: []string{
				"如果 target_status 是 'suspended'，用户不能有状态为 'processing' 的订单。",
				"操作者必须拥有 'admin:users' 权限。",
			},
			SideEffects: []string{
				"目标用户的状态字段将被更改。",
				"系统将通过邮件异步发送状态变更通知给用户。",
			},
			Hints: "在执行 'deleted' 操作前，请务必在你的推理过程中生成一段日志，解释删除的理由。",
		},
		// -----------------------------
		Parameters: []interface{}{
			map[string]interface{}{
				"name":     "id",
				"in":       "path",
				"required": true,
				"schema":   map[string]interface{}{"type": "string"},
			},
		},
		RequestBody: map[string]interface{}{
			"content": map[string]interface{}{
				"application/json": map[string]interface{}{
					"schema": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"status": map[string]interface{}{
								"type": "string",
								"enum": []string{"active", "suspended", "deleted"},
							},
						},
					},
				},
			},
		},
		Responses: map[string]interface{}{
			"200": map[string]interface{}{"description": "Status updated successfully"},
		},
	})
}

func main() {
	r := gin.Default()

	// 暴露 OpenAPI 规范文档的端点
	r.GET("/openapi.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, spec.GlobalSpec.GenerateJSON())
	})

	// 实际的业务处理逻辑 (略)
	// r.PUT("/agentic/v1/users/:id/status", handleUserStatusUpdate)

	log.Println("Agentic OpenAPI Server running on http://localhost:8080")
	log.Println("Fetch the Agent-ready spec at: http://localhost:8080/openapi.json")
	r.Run(":8080")
}
