package dto

// CreateOrderRequest 用于 POST /orders
type CreateOrderRequest struct {
	// 金额通常以"分"为单位，使用 int64
	Amount int64 `json:"amount" binding:"required,gt=0"`
}

// UpdateOrderRequest 用于 PATCH /orders/:id
// 遵循第 02 讲的指针规范，支持零值更新
type UpdateOrderRequest struct {
	Amount *int64 `json:"amount" binding:"omitempty,gt=0"`
	// 注意：Status 字段不在这里暴露
	// 状态流转必须通过专门的 Custom Method (如 Cancel) 进行
}
