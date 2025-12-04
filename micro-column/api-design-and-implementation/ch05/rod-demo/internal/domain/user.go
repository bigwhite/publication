package domain

// 增加一个嵌套结构
type UserProfile struct {
	Avatar string `json:"avatar"`
	City   string `json:"city"`
}

// 升级 User 实体
type User struct {
	ID       string       `json:"id"`
	Bio      string       `json:"bio"`
	Name     string       `json:"name"`
	Age      int          `json:"age"`
	IsActive bool         `json:"is_active"`
	Profile  *UserProfile `json:"profile"` // 嵌套字段
}
