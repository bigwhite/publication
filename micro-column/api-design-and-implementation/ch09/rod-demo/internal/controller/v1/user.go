package v1

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"rod-demo/internal/controller"
	"rod-demo/internal/domain"
	dtov1 "rod-demo/internal/dto/v1"
)

// UserV1Service V1 版本的适配层
// 实现了 controller.CRUDService[dtov1.UserResponse, dtov1.CreateUserRequest, dtov1.UpdateUserRequest]
type UserV1Service struct {
	store map[string]domain.User // 引用共享存储
	mu    *sync.RWMutex
}

// Create: 将 V1 的 Name/Age 转换为 Domain 的 FirstName/LastName/BirthYear
func (s *UserV1Service) Create(ctx *gin.Context, req *dtov1.CreateUserRequest) (*dtov1.UserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()

	// 简单的姓名拆分逻辑 (模拟兼容旧代码)
	names := strings.SplitN(req.Name, " ", 2)
	firstName := names[0]
	lastName := ""
	if len(names) > 1 {
		lastName = names[1]
	}

	// 存入 Domain Model
	user := domain.User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		BirthYear: time.Now().Year() - req.Age, // Age -> BirthYear
	}
	s.store[id] = user

	// 返回 V1 响应
	return &dtov1.UserResponse{
		ID:   user.ID,
		Name: req.Name,
		Age:  req.Age,
	}, nil
}

// Get: 将 Domain 数据组装回 V1 格式
func (s *UserV1Service) Get(ctx *gin.Context, id string) (*dtov1.UserResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.store[id]
	if !exists {
		return nil, nil
	}

	return &dtov1.UserResponse{
		ID:   user.ID,
		Name: user.FirstName + " " + user.LastName, // 拼接
		Age:  time.Now().Year() - user.BirthYear,   // 计算
	}, nil
}

// List: 批量适配
func (s *UserV1Service) List(ctx *gin.Context, req controller.ListRequest) (*controller.ListResponse[*dtov1.UserResponse], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var items []*dtov1.UserResponse
	// 简单全量遍历，生产环境应使用游标分页
	for _, user := range s.store {
		items = append(items, &dtov1.UserResponse{
			ID:   user.ID,
			Name: user.FirstName + " " + user.LastName,
			Age:  time.Now().Year() - user.BirthYear,
		})
	}

	return &controller.ListResponse[*dtov1.UserResponse]{
		Items: items,
		// NextPageToken: ... (省略分页逻辑以聚焦版本化)
	}, nil
}

// Update: 处理 V1 的局部更新
func (s *UserV1Service) Update(ctx *gin.Context, id string, req *dtov1.UpdateUserRequest) (*dtov1.UserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.store[id]
	if !exists {
		return nil, errors.New("not found")
	}

	// 兼容逻辑：如果 V1 传了 Name，需要拆解并更新底层
	if req.Name != nil {
		names := strings.SplitN(*req.Name, " ", 2)
		user.FirstName = names[0]
		if len(names) > 1 {
			user.LastName = names[1]
		}
	}

	// 兼容逻辑：Age -> BirthYear
	if req.Age != nil {
		user.BirthYear = time.Now().Year() - *req.Age
	}

	s.store[id] = user

	return &dtov1.UserResponse{
		ID:   user.ID,
		Name: user.FirstName + " " + user.LastName,
		Age:  time.Now().Year() - user.BirthYear,
	}, nil
}

// Delete: 删除逻辑通常是通用的
func (s *UserV1Service) Delete(ctx *gin.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, id)
	return nil
}

// ---------------------------------------------------------
// V1 控制器：组合 BaseController
// ---------------------------------------------------------

type UserController struct {
	// 泛型参数 T 指定为 dtov1.UserResponse
	*controller.BaseController[dtov1.UserResponse, dtov1.CreateUserRequest, dtov1.UpdateUserRequest]
}

func NewUserController(store map[string]domain.User, mu *sync.RWMutex) *UserController {
	svc := &UserV1Service{store: store, mu: mu}
	// 实例化 BaseController，自动获得 V1 风格的 CRUD HTTP 接口
	base := controller.NewBaseController[dtov1.UserResponse, dtov1.CreateUserRequest, dtov1.UpdateUserRequest](svc)
	return &UserController{BaseController: base}
}
