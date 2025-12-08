package v2

type UserResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"` // 新：拆分
	LastName  string `json:"last_name"`  // 新：拆分
	BirthYear int    `json:"birth_year"` // 新：原始数据
}

type CreateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthYear int    `json:"birth_year"`
}

type UpdateUserRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	BirthYear *int    `json:"birth_year"`
}
