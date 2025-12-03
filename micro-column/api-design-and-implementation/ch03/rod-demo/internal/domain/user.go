package domain

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Bio      string `json:"bio"`
	IsActive bool   `json:"is_active"`
}
