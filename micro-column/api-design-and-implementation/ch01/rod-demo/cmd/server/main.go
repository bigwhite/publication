package main

import (
	"rod-demo/internal/router"
	"rod-demo/internal/user"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// API 版本控制（这也是一种资源层级）
	v1 := r.Group("/api/v1")

	// 注册 Users 资源
	// 只需要这一行，就自动完成了 5 个标准接口的绑定
	router.RegisterResource(v1, "users", user.NewUserController())

	// 如果未来有 Orders 资源，也是一样的模式：
	// router.RegisterResource(v1, "orders", order.NewOrderController())

	r.Run(":8080")
}
