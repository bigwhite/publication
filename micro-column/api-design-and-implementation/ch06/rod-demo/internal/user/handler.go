package user

import (
	"sort"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"rod-demo/internal/controller"
	"rod-demo/internal/domain"
	"rod-demo/internal/dto"
	"rod-demo/pkg/errs"
	"rod-demo/pkg/pagination"
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
	// 为了演示分页，我们需要一个有序的列表作为"数据库索引"
	// 在真实数据库中，这对应 B+ 树索引
	sortedIDs []string
}

func NewUserService() *UserService {
	s := &UserService{
		store:     make(map[string]domain.User),
		sortedIDs: make([]string, 0),
	}
	// 初始化模拟数据：生成 105 个用户，测试分页
	// ID 从 "10001" 到 "10105"
	for i := 0; i < 105; i++ {
		id := strconv.Itoa(10000 + i)
		s.store[id] = domain.User{
			ID:   id,
			Name: "User " + id,
			Age:  18 + (i % 10),
		}
		s.sortedIDs = append(s.sortedIDs, id)
	}
	// 模拟数据库倒序索引 (ORDER BY id DESC)
	// 最新创建的用户排在最前面
	sort.Sort(sort.Reverse(sort.StringSlice(s.sortedIDs)))

	return s
}

// Create 实现具体的创建逻辑
// BaseController 会自动调用此方法，并处理 HTTP 201 响应
func (s *UserService) Create(ctx *gin.Context, req *dto.CreateUserRequest) (*domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// [修改点] 模拟复杂的业务校验
	if req.Age < 0 || req.Age > 150 {
		// 构建结构化的错误详情 (Google AIP 风格)
		details := map[string]interface{}{
			"field":       "age",
			"value":       req.Age,
			"constraint":  "0 <= age <= 150",
			"description": "age value is unrealistic",
		}

		// 返回带详情的错误
		return nil, errs.New(errs.ErrBadRequest, "invalid parameters").
			WithDetails(details)
	}

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
		// [修改点] 不再返回 nil, nil
		// 而是返回带业务语义的结构化错误
		// BaseController 会捕获这个错误并交给中间件处理
		return nil, errs.New(errs.ErrNotFound, "user not found with id "+id)
	}
	return &user, nil
}

// List 实现基于游标的高性能分页逻辑
// 对应 Google AIP-158
func (s *UserService) List(ctx *gin.Context, req controller.ListRequest) (*controller.ListResponse[*domain.User], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 1. 解码 Token 获取游标 (Cursor)
	// 如果 Token 为空，cursor 为空字符串，表示第一页
	cursor, err := pagination.Decode(req.PageToken)
	if err != nil {
		return nil, err // 实际应返回 Invalid Argument 错误
	}

	var result []*domain.User
	var nextCursor string

	// 2. 模拟数据库查询 (Keyset Pagination)
	// SQL 语义: SELECT * FROM users WHERE id < cursor ORDER BY id DESC LIMIT page_size

	// 在 sortedIDs 中找到 cursor 的位置
	startIndex := 0
	if cursor != "" {
		// 二分查找 cursor 的位置 (模拟数据库索引定位)
		// 注意：sortedIDs 是倒序的
		idx := sort.Search(len(s.sortedIDs), func(i int) bool {
			return s.sortedIDs[i] <= cursor
		})
		// sort.Search 返回的是 <= cursor 的第一个位置
		// 如果找到的值等于 cursor，说明是上一页的最后一条，我们需要从它的下一条开始
		if idx < len(s.sortedIDs) && s.sortedIDs[idx] == cursor {
			startIndex = idx + 1
		} else {
			startIndex = idx
		}
	}

	// 3. 截取数据 (LIMIT page_size)
	count := 0
	for i := startIndex; i < len(s.sortedIDs); i++ {
		if count >= req.PageSize {
			break
		}
		id := s.sortedIDs[i]
		if user, ok := s.store[id]; ok {
			// 注意：必须拷贝副本，避免指针指向循环变量
			temp := user
			result = append(result, &temp)
			// 记录当前最后一条数据的 ID，作为下一次的游标
			nextCursor = id
			count++
		}
	}

	// 4. 生成 NextPageToken
	// 如果取出的数据量少于 PageSize，或者已经到了数组末尾，说明没有下一页了
	encodedToken := ""
	if len(result) == req.PageSize && startIndex+len(result) < len(s.sortedIDs) {
		encodedToken = pagination.Encode(nextCursor)
	}

	// 5. 构造标准响应
	return &controller.ListResponse[*domain.User]{
		Items:         result,
		NextPageToken: encodedToken,
		TotalSize:     len(s.sortedIDs), // 可选
	}, nil
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
