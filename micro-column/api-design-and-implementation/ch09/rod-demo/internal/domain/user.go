package domain

// 增加一个嵌套结构
type UserProfile struct {
	Avatar string `json:"avatar"`
	City   string `json:"city"`
}

// 升级 User 实体
// User 是底层共享的领域实体
// 无论 API 版本怎么变，这个结构体代表业务的本质
type User struct {
	ID        string
	FirstName string
	LastName  string
	BirthYear int
	Bio       string       `json:"bio"`
	IsActive  bool         `json:"is_active"`
	Profile   *UserProfile `json:"profile"` // 嵌套字段
}
