package main

import (
	"rod-demo/internal/middleware"
	"rod-demo/internal/router"
	"rod-demo/pkg/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化 Redis
	redis.Init()

	r := gin.Default()

	// 注册全局错误处理中间件
	r.Use(middleware.ErrorHandler())

	// 不再需要在 main 中手动划分 v1/v2 或挂载中间件
	// 所有的版本编排策略，都已封装在 router.RegisterRoutes 中
	router.RegisterRoutes(r)

	r.Run(":8080")
}
