package router

import (
	"rod-demo/internal/controller"

	"github.com/gin-gonic/gin"
)

// RegisterResource 将一个标准控制器注册到 Gin 的路由组中
// resourceName 必须是复数，例如 "users"
func RegisterResource(r *gin.RouterGroup, resourceName string, ctrl controller.StandardController) {
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
}
