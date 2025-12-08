package controller

import (
	"net/http"
	"rod-demo/pkg/errs"
	"rod-demo/pkg/fieldmask"
	"strings"

	"github.com/gin-gonic/gin"
)

// CRUDService 定义了业务层（Service）必须实现的通用接口。
// T: 领域实体类型 (Domain Entity)，例如 domain.User

// 1. 升级 Service 接口：List 方法增加入参和出参
type CRUDService[T any, CreateReq any, UpdateReq any] interface {
	Create(ctx *gin.Context, req *CreateReq) (*T, error)
	Get(ctx *gin.Context, id string) (*T, error)

	// [修改点] List 现在接收 ListRequest 并返回 ListResponse
	List(ctx *gin.Context, req ListRequest) (*ListResponse[*T], error)

	Update(ctx *gin.Context, id string, req *UpdateReq) (*T, error)
	Delete(ctx *gin.Context, id string) error
}

// BaseController 是一个泛型控制器，实现了标准的 RESTful CRUD 操作。
// 具体的 Controller (如 UserController) 可以通过嵌入此结构体来继承标准行为。
type BaseController[T any, CreateReq any, UpdateReq any] struct {
	Service CRUDService[T, CreateReq, UpdateReq]
}

// NewBaseController 创建一个新的泛型控制器实例
func NewBaseController[T any, C any, U any](s CRUDService[T, C, U]) *BaseController[T, C, U] {
	return &BaseController[T, C, U]{
		Service: s,
	}
}

// Create 处理 POST /resources 请求
func (bc *BaseController[T, C, U]) Create(c *gin.Context) {
	var req C
	// 1. 绑定 JSON 请求体到 CreateReq DTO
	if err := c.ShouldBindJSON(&req); err != nil {
		// 关键：必须将原始 error 包装为 AppError，并指定类型为 ErrBadRequest
		// 这样中间件才能识别出这是 400，而不是 500
		c.Error(errs.Wrap(errs.ErrBadRequest, "invalid request parameters", err))
		return
	}

	// 2. 调用 Service 层逻辑
	result, err := bc.Service.Create(c, &req)
	if err != nil {
		// 这里的错误处理比较粗糙，第06讲我们会引入结构化错误处理
		c.Error(err)
		return
	}

	// 3. 成功返回 201 Created 和创建后的资源
	c.JSON(http.StatusCreated, result)
}

// Get 处理 GET /resources/:id 请求
// 升级：支持 ?fields=id,name,profile.city 参数
func (bc *BaseController[T, C, U]) Get(c *gin.Context) {
	id := c.Param("id")

	result, err := bc.Service.Get(c, id)
	if err != nil {
		// 关键：这里将 Service 返回的 AppError 推入 Gin 上下文
		// 之后的 ErrorHandler 中间件会捕获它并渲染响应
		c.Error(err)
		return
	}

	if result == nil {
		c.Error(err)
		return
	}

	// --- 核心修改开始 ---

	// 1. 获取并解析 fields 参数
	// Google AIP-161 建议使用 `read_mask` 或 `fields`
	fieldsParam := c.Query("fields")

	if fieldsParam != "" {
		fields := strings.Split(fieldsParam, ",")

		// 2. 调用工具包进行裁剪
		// 注意：这里 result 是 *T 类型
		prunedData, err := fieldmask.Prune(result, fields)
		if err != nil {
			// 裁剪失败（通常是 JSON 序列化问题），降级返回完整数据或报错
			// 这里选择报错以便调试
			c.Error(err)
			return
		}

		// 3. 返回裁剪后的 Map
		c.JSON(http.StatusOK, prunedData)
		return
	}

	// --- 核心修改结束 ---

	// 如果没有 fields 参数，默认返回全量数据
	c.JSON(http.StatusOK, result)
}

// List 处理 GET /resources 请求
func (bc *BaseController[T, C, U]) List(c *gin.Context) {

	var req ListRequest

	// 1. 绑定 Query 参数 (page_size, page_token)
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认 PageSize
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// 调用 Service 层获取列表
	results, err := bc.Service.List(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 成功返回 200 OK
	// 即使列表为空，也应返回空数组 []，而不是 null
	if results == nil {
		c.JSON(http.StatusOK, []T{})
		return
	}

	// --- 核心修改开始 ---
	fieldsParam := c.Query("fields")
	if fieldsParam != "" {
		fields := strings.Split(fieldsParam, ",")
		// Prune 函数同时支持 Struct 和 Slice，直接传入即可
		prunedData, err := fieldmask.Prune(results, fields)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, prunedData)
		return
	}
	// --- 核心修改结束 ---

	c.JSON(http.StatusOK, results)
}

// Update 处理 PATCH /resources/:id 请求
func (bc *BaseController[T, C, U]) Update(c *gin.Context) {
	id := c.Param("id")

	var req U
	// [修改点 2] 同理，Update 也要适配
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.Wrap(errs.ErrBadRequest, "invalid request parameters", err))
		return
	}

	// 2. 调用 Service 执行局部更新
	result, err := bc.Service.Update(c, id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. 成功返回 200 OK 和更新后的完整资源
	c.JSON(http.StatusOK, result)
}

// Delete 处理 DELETE /resources/:id 请求
func (bc *BaseController[T, C, U]) Delete(c *gin.Context) {
	id := c.Param("id")

	// 1. 调用 Service 执行删除
	err := bc.Service.Delete(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 2. 成功返回 204 No Content
	// 根据 HTTP 语义，删除成功后不需要返回 Body
	c.Status(http.StatusNoContent)
}
