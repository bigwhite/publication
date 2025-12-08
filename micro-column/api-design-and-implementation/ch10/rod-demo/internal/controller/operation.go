package controller

import (
	"rod-demo/internal/domain"
	"rod-demo/internal/task"
	"rod-demo/pkg/errs"

	"github.com/gin-gonic/gin"
)

// OperationService 适配器
// 泛型参数：Domain=Operation, CreateReq=any, UpdateReq=any
type OperationService struct{}

// Get 实现标准的查询逻辑
func (s *OperationService) Get(ctx *gin.Context, id string) (*domain.Operation, error) {
	op := task.GetOperation(id)
	if op == nil {
		return nil, errs.New(errs.ErrNotFound, "operation not found")
	}
	return op, nil
}

// -------------------------------------------------------------------
// 下面的方法对于 Operation 资源来说不支持，返回 Method Not Allowed 即可
// -------------------------------------------------------------------

func (s *OperationService) Create(ctx *gin.Context, req *any) (*domain.Operation, error) {
	return nil, errs.New(errs.ErrBadRequest, "cannot create operation directly")
}
func (s *OperationService) List(ctx *gin.Context, req ListRequest) (*ListResponse[*domain.Operation], error) {
	return nil, nil // 简化处理，实际项目可支持列表查询
}
func (s *OperationService) Update(ctx *gin.Context, id string, req *any) (*domain.Operation, error) {
	return nil, errs.New(errs.ErrBadRequest, "cannot update operation")
}
func (s *OperationService) Delete(ctx *gin.Context, id string) error {
	return errs.New(errs.ErrBadRequest, "cannot delete operation")
}

// OperationController
// 继承 BaseController，自动获得标准的 GET /operations/:id 能力
type OperationController struct {
	*BaseController[domain.Operation, any, any]
}

func NewOperationController() *OperationController {
	svc := &OperationService{}
	base := NewBaseController[domain.Operation, any, any](svc)
	return &OperationController{BaseController: base}
}
