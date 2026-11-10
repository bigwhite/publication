package orders

import "time"

type Order struct {
	ID          string
	CustomerID  string
	Status      string
	TotalAmount int64
	Currency    string
	CreatedAt   time.Time
}
