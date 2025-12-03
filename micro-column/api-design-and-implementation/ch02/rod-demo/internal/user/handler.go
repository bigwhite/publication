package user

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"rod-demo/internal/controller"
	"rod-demo/internal/domain"
	"rod-demo/internal/dto"
)

// =============================================================================
// 1. 定义业务服务 (UserService)
//    它实现了 controller.CRUDService 接口，负责具体的业务逻辑
// =============================================================================

type UserService struct {
	// 模拟数据库存储
	// 在实际项目中，这里通常会注入 Repository 层
	mu    sync.RWMutex
	store map[string]domain.User
}

func NewUserService() *UserService {
	return &UserService{
		store: map[string]domain.User{
			"123": {ID: "123", Name: "Old Tony", Age: 18, Bio: "Original", IsActive: true},
		},
	}
}

// Create 实现具体的创建逻辑
// BaseController 会自动调用此方法，并处理 HTTP 201 响应
func (s *UserService) Create(ctx *gin.Context, req *dto.CreateUserRequest) (*domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 模拟生成 ID
	id := uuid.New().String()

	user := domain.User{
		ID:       id,
		Name:     req.Name,
		Age:      req.Age,
		Bio:      "Default Bio",
		IsActive: true,
	}

	// 存入模拟数据库
	s.store[id] = user
	return &user, nil
}

// Get 实现具体的查询逻辑
// BaseController 会自动处理 HTTP 200 或 404
func (s *UserService) Get(ctx *gin.Context, id string) (*domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.store[id]
	if !exists {
		// 返回 nil, nil 会被 BaseController 识别为 404 Not Found
		return nil, nil
	}
	return &user, nil
}

// List 实现具体的列表查询逻辑
func (s *UserService) List(ctx *gin.Context) ([]*domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*domain.User, 0, len(s.store))
	for _, u := range s.store {
		// 注意：循环变量地址问题，需拷贝
		temp := u
		users = append(users, &temp)
	}
	return users, nil
}

// Update 实现具体的更新逻辑 (核心是处理 Pointer DTO)
// BaseController 会自动绑定 JSON 到 req，并处理错误
func (s *UserService) Update(ctx *gin.Context, id string, req *dto.UpdateUserRequest) (*domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. 检查资源是否存在
	user, exists := s.store[id]
	if !exists {
		return nil, nil // 触发 404
	}

	// 2. 核心逻辑：零值更新处理
	// 只有当 DTO 中的指针不为 nil 时，才更新对应的字段
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Age != nil {
		// 即使是 0 也能正确更新，因为指针非 nil
		user.Age = *req.Age
	}
	if req.Bio != nil {
		user.Bio = *req.Bio
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// 3. 保存回数据库
	s.store[id] = user

	return &user, nil
}

// Delete 实现具体的删除逻辑
func (s *UserService) Delete(ctx *gin.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 即使 ID 不存在，Delete 通常也视为成功（幂等性）
	delete(s.store, id)
	return nil
}

// =============================================================================
// 2. 定义控制器 (UserController)
//    通过组合泛型 BaseController，自动获得标准 HTTP 处理能力
// =============================================================================

type UserController struct {
	// 核心：嵌入泛型 BaseController
	// 这行代码让 UserController 直接拥有了 Create, Get, List, Update, Delete 五个 HTTP Handler
	*controller.BaseController[domain.User, dto.CreateUserRequest, dto.UpdateUserRequest]
}

func NewUserController() *UserController {
	// 1. 实例化业务服务
	userService := NewUserService()

	// 2. 实例化泛型基础控制器，注入业务服务
	base := controller.NewBaseController[domain.User, dto.CreateUserRequest, dto.UpdateUserRequest](userService)

	// 3. 返回组合后的控制器
	return &UserController{
		BaseController: base,
	}
}
