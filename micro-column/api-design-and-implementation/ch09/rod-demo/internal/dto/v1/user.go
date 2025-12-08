package v1

type UserResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"` // 旧：合并的姓名
	Age  int    `json:"age"`  // 旧：动态计算的年龄
}

type CreateUserRequest struct {
	Name string `json:"name"` // 输入全名
	Age  int    `json:"age"`
}

type UpdateUserRequest struct {
	Name *string `json:"name"`
	Age  *int    `json:"age"`
}
