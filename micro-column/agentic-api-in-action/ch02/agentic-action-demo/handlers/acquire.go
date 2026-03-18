package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FetchOrderRequest 定义了获取订单所需的上下文
type FetchOrderRequest struct {
	OrderID   string `json:"order_id" binding:"required"`
	WithItems bool   `json:"with_items"` // AI 可以决定是否需要明细
}

// FetchOrder 是一个典型的数据获取 Action
func FetchOrder(c *gin.Context) {
	var req FetchOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid context for FETCH action: " + err.Error()})
		return
	}

	// 模拟数据库查询...
	orderData := gin.H{
		"order_id": req.OrderID,
		"status":   "SHIPPED",
		"amount":   199.50,
	}

	if req.WithItems {
		orderData["items"] = []string{"Mechanical Keyboard", "Mousepad"}
	}

	// 返回结构化的意图响应
	c.JSON(http.StatusOK, gin.H{
		"action": "FETCH",
		"status": "SUCCESS",
		"data":   orderData,
	})
}

// SearchOrders 模拟搜索操作
func SearchOrders(c *gin.Context) {
	// 实现略...
	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "data": "List of orders..."})
}
