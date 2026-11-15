package cancellations

// Request 表示一条订单取消申请。
// 注意：创建 Request 不代表订单已经取消。
type Request struct {
	ID             string `json:"id"`
	OrderID        string `json:"order_id"`
	Reason         string `json:"reason"`
	Status         string `json:"status"`
	IdempotencyKey string `json:"-"`
}
