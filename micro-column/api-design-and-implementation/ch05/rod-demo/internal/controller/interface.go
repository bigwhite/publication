package controller

import "github.com/gin-gonic/gin"

// StandardController 定义了资源导向设计的 5 个标准方法
// 任何资源控制器都应当实现该接口
type StandardController interface {
	// List 列出集合中的资源 (GET /collection)
	List(c *gin.Context)

	// Get 获取单个资源 (GET /collection/:id)
	Get(c *gin.Context)

	// Create 新建资源 (POST /collection)
	Create(c *gin.Context)

	// Update 更新资源 (PATCH /collection/:id)
	Update(c *gin.Context)

	// Delete 删除资源 (DELETE /collection/:id)
	Delete(c *gin.Context)
}
