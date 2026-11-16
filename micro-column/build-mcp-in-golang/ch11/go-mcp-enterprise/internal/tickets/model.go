package tickets

import "time"

type Ticket struct {
	ID          string
	OrderID     string
	Title       string
	Description string
	Priority    string
	Status      string
	CreatedAt   time.Time
}

type CreateInput struct {
	OrderID        string
	Title          string
	Description    string
	Priority       string
	IdempotencyKey string
}
