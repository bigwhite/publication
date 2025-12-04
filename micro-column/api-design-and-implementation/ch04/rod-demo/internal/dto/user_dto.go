package dto

// CreateUserRequest 用于 POST 请求，字段通常是必填的
type CreateUserRequest struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age" binding:"gte=0,lte=150"`
}

// UpdateUserRequest 用于 PATCH 请求，所有字段均为指针，且可选
type UpdateUserRequest struct {
	Name     *string `json:"name"`      // nil: 不更新; non-nil: 更新
	Age      *int    `json:"age"`       // nil: 不更新; 0: 更新为0
	Bio      *string `json:"bio"`       // 用于演示更新为空字符串
	IsActive *bool   `json:"is_active"` // 用于演示 bool 值的更新
}
