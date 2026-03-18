package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RefundOrderRequest 定义了退款交易所需的严谨参数
type RefundOrderRequest struct {
	OrderID string  `json:"order_id" binding:"required"`
	Amount  float64 `json:"amount" binding:"required,gt=0"`
	Reason  string  `json:"reason" binding:"required"`
}

// RefundOrder 是一个典型的 TRANSACT Action，会改变系统状态
func RefundOrder(c *gin.Context) {
	var req RefundOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 返回详细的验证错误，帮助 AI 进行自我修正 (Self-Correction)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "TRANSACT action failed due to invalid parameters: " + err.Error(),
			"hint":  "Ensure order_id, amount (>0), and reason are provided.",
		})
		return
	}

	// 在这里，我们可以轻松地添加针对 TRANSACT 类别的高级审计日志或二次确认机制
	// 模拟退款业务逻辑...

	c.JSON(http.StatusOK, gin.H{
		"action":  "REFUND",
		"status":  "SUCCESS",
		"message": "Refund processed successfully for order: " + req.OrderID,
		// 返回副作用的结果，让 AI 明确知道操作成功了
		"side_effects": gin.H{
			"refunded_amount":  req.Amount,
			"new_order_status": "REFUNDED",
		},
	})
}
