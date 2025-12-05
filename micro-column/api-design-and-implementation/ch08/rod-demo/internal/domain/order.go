package domain

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusPaid      OrderStatus = "PAID"
	OrderStatusCancelled OrderStatus = "CANCELLED"
	OrderStatusShipped   OrderStatus = "SHIPPED"
)

type Order struct {
	ID     string      `json:"id"`
	Amount int64       `json:"amount"`
	Status OrderStatus `json:"status"`
}

// CanCancel 判断订单是否可以取消
// 将状态流转逻辑封装在领域模型中，而不是散落在 Controller 里
func (o *Order) CanCancel() bool {
	return o.Status == OrderStatusPending || o.Status == OrderStatusPaid
}
