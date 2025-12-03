package order

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"rod-demo/internal/controller"
	"rod-demo/internal/domain"
	"rod-demo/internal/dto"
)

// =============================================================================
// 1. 业务服务 (OrderService)
//    实现 controller.CRUDService 接口，同时包含自定义的 CancelOrder 逻辑
// =============================================================================

type OrderService struct {
	mu    sync.RWMutex
	store map[string]domain.Order
}

func NewOrderService() *OrderService {
	// 初始化一些模拟数据，方便演示
	return &OrderService{
		store: map[string]domain.Order{
			"1001": {ID: "1001", Amount: 9900, Status: domain.OrderStatusPending},
			"1002": {ID: "1002", Amount: 100, Status: domain.OrderStatusShipped},
		},
	}
}

// --- 标准 CRUD 实现 (满足 BaseController 泛型约束) ---

func (s *OrderService) Create(ctx *gin.Context, req *dto.CreateOrderRequest) (*domain.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	order := domain.Order{
		ID:     id,
		Amount: req.Amount,
		Status: domain.OrderStatusPending, // 默认状态
	}
	s.store[id] = order
	return &order, nil
}

func (s *OrderService) Get(ctx *gin.Context, id string) (*domain.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, exists := s.store[id]
	if !exists {
		return nil, nil // BaseController 将转为 404
	}
	return &order, nil
}

func (s *OrderService) List(ctx *gin.Context) ([]*domain.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	orders := make([]*domain.Order, 0, len(s.store))
	for _, o := range s.store {
		temp := o
		orders = append(orders, &temp)
	}
	return orders, nil
}

func (s *OrderService) Update(ctx *gin.Context, id string, req *dto.UpdateOrderRequest) (*domain.Order, error) {
	// 这里的 Update 通常用于修改金额或备注，但不应该用于流转状态
	// 状态流转应通过自定义方法（如 Cancel）进行
	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.store[id]
	if !exists {
		return nil, nil
	}

	if req.Amount != nil {
		order.Amount = *req.Amount
	}

	s.store[id] = order
	return &order, nil
}

func (s *OrderService) Delete(ctx *gin.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, id)
	return nil
}

// --- 自定义业务逻辑 (非 CRUD) ---

// CancelOrder 执行取消逻辑
// 返回值：更新后的订单，以及特定的业务错误（如状态冲突）
func (s *OrderService) CancelOrder(id string) (*domain.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.store[id]
	if !exists {
		return nil, errors.New("not_found")
	}

	// 领域逻辑：检查状态机
	if !order.CanCancel() {
		return nil, errors.New("conflict")
	}

	// 执行副作用（如退款、库存释放），此处省略...

	// 更新状态
	order.Status = domain.OrderStatusCancelled
	s.store[id] = order

	return &order, nil
}

// =============================================================================
// 2. 控制器 (OrderController)
//    继承 BaseController 的能力，并扩展自定义方法 Cancel
// =============================================================================

type OrderController struct {
	// 组合泛型 BaseController，自动拥有 standard methods
	*controller.BaseController[domain.Order, dto.CreateOrderRequest, dto.UpdateOrderRequest]

	// 额外持有具体的 service 引用，以便调用 CancelOrder 等非 CRUD 方法
	service *OrderService
}

func NewOrderController() *OrderController {
	svc := NewOrderService()
	base := controller.NewBaseController[domain.Order, dto.CreateOrderRequest, dto.UpdateOrderRequest](svc)

	return &OrderController{
		BaseController: base,
		service:        svc,
	}
}

// Cancel 处理取消订单的自定义方法
// 映射路由: POST /orders/:id/cancel
func (ctrl *OrderController) Cancel(c *gin.Context) {
	id := c.Param("id")

	// 调用 Service 层的自定义逻辑
	order, err := ctrl.service.CancelOrder(id)
	if err != nil {
		if err.Error() == "not_found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
		if err.Error() == "conflict" {
			// 违反业务规则，返回 409 Conflict
			c.JSON(http.StatusConflict, gin.H{
				"error": "order cannot be cancelled in current status",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回更新后的资源
	c.JSON(http.StatusOK, order)
}
