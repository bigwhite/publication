package domain

import (
	"time"
)

type Link struct {
	ID          int64
	OriginalURL string
	ShortCode   string
	CreatedAt   time.Time
}
