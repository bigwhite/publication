package main

import (
	"agentic-gateway/adapter"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 这是专门给 AI 开设的 VIP 通道
	agenticAPI := r.Group("/agentic/v1")
	{
		// 我们不再暴露危险的 PATCH /users，而是暴露意图明确的 SUSPEND 动作
		agenticAPI.POST("/transact/suspend_user", adapter.SuspendUserAction)
	}

	log.Println("Agentic Gateway Server running on http://localhost:8080")
	r.Run(":8080")
}
