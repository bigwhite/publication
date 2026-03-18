package main

import (
	"agentic-dryrun-demo/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	agenticAPI := r.Group("/agentic/v1")
	{
		// 注册高危的支付动作接口
		agenticAPI.POST("/transact/payment", handlers.ProcessPayment)
	}

	log.Println("Agentic Dry Run Server running on http://localhost:8080")
	r.Run(":8080")
}
