package adapter

import (
	"agentic-gateway/legacy"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuspendUserRequest 是暴露给 AI 的、意图极其清晰的强类型契约
type SuspendUserRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Reason string `json:"reason" binding:"required"` // 强制 AI 必须提供封号理由
}

// SuspendUserAction 充当了“翻译官”和“保镖”
func SuspendUserAction(c *gin.Context) {
	var req SuspendUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse SUSPEND intent. Missing user_id or reason.",
		})
		return
	}

	// 1. 业务护栏 (Guardrails)
	// 比如：禁止封禁 ID 为 "admin_001" 的超级管理员
	if req.UserID == "admin_001" {
		c.JSON(http.StatusForbidden, gin.H{
			"status": "FAILED",
			"error":  "Action DENIED. Cannot suspend a super admin account.",
		})
		return
	}

	// 2. 意图翻译 (Intent Translation)
	// 我们在这里，安全地、确定性地将明确的意图，翻译为遗留系统需要的底层数据结构
	falseVal := false
	legacyPayload := legacy.UserUpdatePayload{
		IsActive: &falseVal,
		// 我们坚决不给 Name 和 Email 赋值，彻底杜绝 AI 误修改其他字段的可能
	}

	// 3. 记录审计日志：AI 是因为什么原因 (Reason) 触发了这个底层调用的
	fmt.Printf("[AUDIT] AI requested SUSPEND for User: %s, Reason: %s\n", req.UserID, req.Reason)

	// 4. 调用遗留系统
	err := legacy.PatchUserStatus(req.UserID, legacyPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "FAILED",
			"error":  "Legacy system failed: " + err.Error(),
		})
		return
	}

	// 5. 返回 Agentic 标准响应
	c.JSON(http.StatusOK, gin.H{
		"action":  "SUSPEND",
		"status":  "SUCCESS",
		"message": "User suspended successfully via legacy adapter.",
	})
}
