package router

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"rod-demo/internal/controller"
	v1 "rod-demo/internal/controller/v1"
	v2 "rod-demo/internal/controller/v2"
	"rod-demo/internal/domain"
	"rod-demo/internal/middleware"
)

// =============================================================================
// 1. 通用路由注册工具 (Infrastructure Layer)
//    这段代码是框架级的，写一次，整个项目生命周期复用。
//    它确保了 V1 和 V2 遵循完全相同的 RESTful URL 规范。
// =============================================================================

// Option 定义路由注册的选项函数 (支持自定义方法扩展，如 03 讲所述)
type Option func(*gin.RouterGroup)

// WithCustomMethod 用于注册自定义方法 (e.g., POST /:id/cancel)
func WithCustomMethod(method, action string, handler gin.HandlerFunc) Option {
	return func(g *gin.RouterGroup) {
		g.Handle(method, "/:id/"+action, handler)
	}
}

// RegisterResource 将一个标准控制器注册到 Gin 的路由组中
// 核心复用逻辑：无论 V1 还是 V2，只要实现了 StandardController，就能自动注册
func RegisterResource(r *gin.RouterGroup, resourceName string, ctrl controller.StandardController, options ...Option) {
	group := r.Group("/" + resourceName)
	{
		// 集合操作
		group.GET("", ctrl.List)
		group.POST("", ctrl.Create)

		// 单个资源操作
		group.GET("/:id", ctrl.Get)
		group.PATCH("/:id", ctrl.Update)
		group.DELETE("/:id", ctrl.Delete)
	}

	// 应用自定义选项 (如 Cancel, Ship 等)
	for _, opt := range options {
		opt(group)
	}
}

// =============================================================================
// 2. 业务路由编排 (Application Layer)
//    这里是版本策略的控制中心，负责组装 V1 和 V2。
// =============================================================================

func RegisterRoutes(r *gin.Engine) {
	// 1. 初始化基础设施 (模拟共享数据库)
	// 在真实场景中，这里通常是 *gorm.DB 或 *redis.Client
	sharedStore := make(map[string]domain.User)
	var mu sync.RWMutex

	// 初始化一条测试数据，确保 ID "1" 存在
	sharedStore["1"] = domain.User{
		ID:        "1",
		FirstName: "Tony",
		LastName:  "Bai",
		BirthYear: 1988,
	}

	// 2. 配置废弃策略 (Aggressive Obsolescence)
	// 设定 V1 将于 2025年底下线
	sunsetDate := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
	deprecationMiddleware := middleware.Deprecator(middleware.DeprecationOptions{
		SunsetDate: sunsetDate,
		Link:       "/api/v2/users",
	})

	// 3. 注册 V1 版本 (Deprecated)
	// 挂载废弃中间件
	v1Group := r.Group("/api/v1", deprecationMiddleware)
	{
		// 注入共享存储，实例化 V1 控制器
		userCtrlV1 := v1.NewUserController(sharedStore, &mu)

		// 复用注册器！
		RegisterResource(v1Group, "users", userCtrlV1)
	}

	// 4. 注册 V2 版本 (Current)
	v2Group := r.Group("/api/v2")
	{
		// 注入同一个共享存储，实例化 V2 控制器
		userCtrlV2 := v2.NewUserController(sharedStore, &mu)

		// 复用注册器！
		RegisterResource(v2Group, "users", userCtrlV2)
	}
}
