package v2

import (
	"errors"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"rod-demo/internal/controller"
	"rod-demo/internal/domain"
	dtov2 "rod-demo/internal/dto/v2"
)

// UserV2Service V2 版本的业务服务
// 实现了 controller.CRUDService[dtov2.UserResponse, dtov2.CreateUserRequest, dtov2.UpdateUserRequest]
type UserV2Service struct {
	store map[string]domain.User // 引用同一个共享存储
	mu    *sync.RWMutex
}

// Create: V2 逻辑，直接映射
func (s *UserV2Service) Create(ctx *gin.Context, req *dtov2.CreateUserRequest) (*dtov2.UserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	user := domain.User{
		ID:        id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		BirthYear: req.BirthYear, // 直接存储年份
	}
	s.store[id] = user

	return &dtov2.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		BirthYear: user.BirthYear,
	}, nil
}

// Get: V2 逻辑，直接映射
func (s *UserV2Service) Get(ctx *gin.Context, id string) (*dtov2.UserResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.store[id]
	if !exists {
		return nil, nil
	}

	return &dtov2.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		BirthYear: user.BirthYear,
	}, nil
}

// List: V2 列表
func (s *UserV2Service) List(ctx *gin.Context, req controller.ListRequest) (*controller.ListResponse[*dtov2.UserResponse], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var items []*dtov2.UserResponse
	for _, user := range s.store {
		items = append(items, &dtov2.UserResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			BirthYear: user.BirthYear,
		})
	}

	return &controller.ListResponse[*dtov2.UserResponse]{
		Items: items,
	}, nil
}

// Update: V2 局部更新
func (s *UserV2Service) Update(ctx *gin.Context, id string, req *dtov2.UpdateUserRequest) (*dtov2.UserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.store[id]
	if !exists {
		return nil, errors.New("not found")
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.BirthYear != nil {
		user.BirthYear = *req.BirthYear
	}

	s.store[id] = user

	return &dtov2.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		BirthYear: user.BirthYear,
	}, nil
}

// Delete: 删除逻辑
func (s *UserV2Service) Delete(ctx *gin.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, id)
	return nil
}

// ---------------------------------------------------------
// V2 控制器：组合 BaseController
// ---------------------------------------------------------

type UserController struct {
	// 泛型参数 T 指定为 dtov2.UserResponse
	*controller.BaseController[dtov2.UserResponse, dtov2.CreateUserRequest, dtov2.UpdateUserRequest]
}

func NewUserController(store map[string]domain.User, mu *sync.RWMutex) *UserController {
	svc := &UserV2Service{store: store, mu: mu}
	// 实例化 BaseController，自动获得 V2 风格的 CRUD HTTP 接口
	base := controller.NewBaseController[dtov2.UserResponse, dtov2.CreateUserRequest, dtov2.UpdateUserRequest](svc)
	return &UserController{BaseController: base}
}
