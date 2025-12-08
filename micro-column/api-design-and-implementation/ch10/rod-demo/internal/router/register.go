package router

import (
	"rod-demo/internal/controller"
	"rod-demo/internal/user"

	"github.com/gin-gonic/gin"
)

// Option 定义路由注册的选项函数
type Option func(*gin.RouterGroup)

// WithCustomMethod 用于注册自定义方法
// method: HTTP动词 (e.g., "POST")
// action: 动作名称 (e.g., "cancel")
// handler: 处理函数
func WithCustomMethod(method, action string, handler gin.HandlerFunc) Option {
	return func(g *gin.RouterGroup) {
		// 注册路径: /:id/action
		// 例如: POST /orders/:id/cancel
		g.Handle(method, "/:id/"+action, handler)
	}
}

// RegisterResource 将一个标准控制器注册到 Gin 的路由组中
// resourceName 必须是复数，例如 "users"
func RegisterResource(r *gin.RouterGroup, resourceName string, ctrl controller.StandardController, options ...Option) {
	// 创建资源集合的路由组，例如 /api/v1/users
	// 这里体现了 ROD 的层级思想：URL 即资源路径
	group := r.Group("/" + resourceName)
	{
		// 集合操作
		group.GET("", ctrl.List)    // GET /users
		group.POST("", ctrl.Create) // POST /users

		// 单个资源操作，:id 代表资源标识符
		group.GET("/:id", ctrl.Get)       // GET /users/:id
		group.PATCH("/:id", ctrl.Update)  // PATCH /users/:id
		group.DELETE("/:id", ctrl.Delete) // DELETE /users/:id
	}

	// 应用自定义选项
	for _, opt := range options {
		opt(group)
	}
}

// SetupRoutes 路由注册入口
func SetupRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")

	// 1. 注册 User 资源 (第 08 讲内容)
	userCtrl := user.NewUserController()
	RegisterResource(v1, "users", userCtrl)

	// --- 第 10 讲新增：LRO 路由 ---

	// 2. 实例化控制器
	opCtrl := controller.NewOperationController()
	aiCtrl := controller.NewAIController()

	// 3. 注册 Operations 资源 (复用 BaseController 能力)
	// 这会自动生成 GET /api/v1/operations/:id
	RegisterResource(v1, "operations", opCtrl)

	// 4. 注册 AI 自定义方法 (LRO 触发)
	// 对应 Google AIP 风格: POST /images:generate
	v1.POST("/images/generate", aiCtrl.GenerateImage)
	v1.POST("/chat/stream", aiCtrl.ChatStream)
}
