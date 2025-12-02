package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserController 实现 controller.StandardController 接口
type UserController struct {
	// 这里可以注入 Service 或 Repository
}

func NewUserController() *UserController {
	return &UserController{}
}

// List: GET /users
func (ctrl *UserController) List(c *gin.Context) {
	// 模拟返回用户列表
	c.JSON(http.StatusOK, gin.H{
		"data":  []string{"Tony Bai", "Gopher"},
		"total": 2,
	})
}

// Get: GET /users/:id
func (ctrl *UserController) Get(c *gin.Context) {
	id := c.Param("id")
	// 资源导向设计中，ID 是资源的核心标识
	c.JSON(http.StatusOK, gin.H{
		"id":   id,
		"name": "Tony Bai",
		"role": "Architect",
	})
}

// Create: POST /users
func (ctrl *UserController) Create(c *gin.Context) {
	// 标准 Create 操作成功后返回 201 Created
	c.JSON(http.StatusCreated, gin.H{
		"status": "created",
		"id":     "1001", // 假设新创建的 ID
	})
}

// Update: PATCH /users/:id
func (ctrl *UserController) Update(c *gin.Context) {
	id := c.Param("id")
	// PATCH 语义：局部更新
	c.JSON(http.StatusOK, gin.H{
		"id":     id,
		"status": "updated",
	})
}

// Delete: DELETE /users/:id
func (ctrl *UserController) Delete(c *gin.Context) {
	// 标准 Delete 操作通常返回 204 No Content (无响应体)
	// 但为了演示效果，这里返回 200
	c.JSON(http.StatusOK, gin.H{
		"status": "deleted",
		"id":     c.Param("id"),
	})
}
