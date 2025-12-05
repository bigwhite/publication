package main

import (
	"rod-demo/internal/middleware"
	"rod-demo/internal/order"
	"rod-demo/internal/router"
	"rod-demo/internal/user"
	"rod-demo/pkg/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化 Redis
	redis.Init()

	r := gin.Default()

	// 注册全局错误处理中间件
	r.Use(middleware.ErrorHandler())

	// API 版本控制（这也是一种资源层级）
	v1 := r.Group("/api/v1")

	// 使用幂等性中间件
	// 通常只针对非 GET/DELETE 请求启用
	v1.Use(middleware.Idempotency())

	// 注册 Users 资源
	// 只需要这一行，就自动完成了 5 个标准接口的绑定
	router.RegisterResource(v1, "users", user.NewUserController())

	orderCtrl := order.NewOrderController()

	// 注册资源，并附带自定义方法
	router.RegisterResource(v1, "orders", orderCtrl,
		// 显式注册 Cancel 动作
		router.WithCustomMethod("POST", "cancel", orderCtrl.Cancel),

		// 如果有其他动作，比如 "发货"
		// router.WithCustomMethod("POST", "ship", orderCtrl.Ship),
	)

	r.Run(":8080")
}
