package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PaymentRequest 定义了支付所需的参数，以及极其重要的 TestMode 标志
type PaymentRequest struct {
	UserID string  `json:"user_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
	// TestMode 如果为 true，则仅模拟执行，不产生真实扣款
	TestMode bool `json:"test_mode"`
}

// 模拟的后端数据库状态
var userBalance = map[string]float64{
	"user_123": 100.00,
}

// ProcessPayment 处理支付请求，支持 Dry Run 模式
func ProcessPayment(c *gin.Context) {
	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment parameters"})
		return
	}

	// 1. 第一阶段：无论是否是 TestMode，都必须执行的【业务校验逻辑】

	// 校验 1：用户是否存在
	currentBalance, exists := userBalance[req.UserID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "FAILED",
			"error":  fmt.Sprintf("User %s not found", req.UserID),
		})
		return
	}

	// 校验 2：余额是否充足
	if currentBalance < req.Amount {
		c.JSON(http.StatusPaymentRequired, gin.H{
			"status": "FAILED",
			"error":  "Insufficient balance",
			"details": gin.H{
				"current_balance": currentBalance,
				"shortfall":       req.Amount - currentBalance,
			},
		})
		return
	}

	// 2. 第二阶段：Dry Run 拦截点 (The Safety Net)

	// 如果是测试模式，到此为止！构造预测报告并返回。
	if req.TestMode {
		predictedNewBalance := currentBalance - req.Amount

		// 模拟生成一些警告信息，供 AI 决策参考
		var warnings []string
		if predictedNewBalance < 10.0 {
			warnings = append(warnings, "Warning: Balance will drop below $10.00 after this transaction.")
		}

		c.JSON(http.StatusOK, gin.H{
			"action":  "TRANSACT",
			"status":  "SIMULATION_SUCCESS",
			"message": "Dry run completed successfully. No funds were actually deducted.",
			"predicted_side_effects": gin.H{
				"deducted_amount":       req.Amount,
				"predicted_new_balance": predictedNewBalance,
			},
			"warnings": warnings,
		})
		return
	}

	// 3. 第三阶段：只有真实的生产请求，才会执行【不可逆的副作用】(The Side Effects)

	// 实际扣款 (模拟更新数据库)
	userBalance[req.UserID] = currentBalance - req.Amount
	actualNewBalance := userBalance[req.UserID]

	// 记录审计日志 (在真实系统中极其重要)
	fmt.Printf("[AUDIT] Real transaction executed for User: %s, Amount: $%.2f\n", req.UserID, req.Amount)

	c.JSON(http.StatusOK, gin.H{
		"action":  "TRANSACT",
		"status":  "SUCCESS",
		"message": "Payment processed successfully.",
		"result": gin.H{
			"deducted_amount": req.Amount,
			"new_balance":     actualNewBalance,
		},
	})
}
