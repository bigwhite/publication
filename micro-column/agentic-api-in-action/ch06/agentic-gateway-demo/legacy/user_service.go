package legacy

import (
	"fmt"
	"log"
)

// 这是一个极其底层的数据库模型映射
type UserUpdatePayload struct {
	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"` // 危险：如果传 false 就封号
}

// 模拟遗留系统的内部接口，它对 AI 极度不友好
func PatchUserStatus(userID string, payload UserUpdatePayload) error {
	log.Printf("[LEGACY SYSTEM] PATCH /legacy/api/users/%s called with payload: %+v\n", userID, payload)

	// 如果不小心把 email 置空了... 灾难！
	if payload.Email != nil && *payload.Email == "" {
		return fmt.Errorf("FATAL: Email cannot be empty string")
	}

	if payload.IsActive != nil {
		status := "Active"
		if !*payload.IsActive {
			status = "Suspended"
		}
		log.Printf("[LEGACY SYSTEM] User %s status changed to: %s\n", userID, status)
	}

	return nil
}
