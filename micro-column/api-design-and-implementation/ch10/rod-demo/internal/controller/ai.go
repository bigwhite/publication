package controller

import (
	"io"
	"net/http"
	"rod-demo/internal/task"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AIController struct{}

func NewAIController() *AIController {
	return &AIController{}
}

// GenerateImage 处理 LRO 任务提交
// 映射路由: POST /images:generate
func (c *AIController) GenerateImage(ctx *gin.Context) {
	// 1. 模拟参数绑定 (实际应定义 DTO)
	type ImageReq struct {
		Prompt string `json:"prompt"`
	}
	var req ImageReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. 生成 Operation ID
	opID := uuid.New().String()

	// 3. 触发后台异步任务
	task.StartImageGeneration(opID, req.Prompt)

	// 4. 立即返回 LRO 对象 (初始状态)
	op := task.GetOperation(opID)
	ctx.JSON(http.StatusOK, op)
}

// ChatStream 处理流式对话
// Custom Method: POST /chat:stream
func (c *AIController) ChatStream(ctx *gin.Context) {
	// 1. 设置 SSE 专用 Header
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")
	// 告诉 Nginx 不要在缓冲区等待，立即下发
	ctx.Header("X-Accel-Buffering", "no")

	// 2. 开启流
	// c.Stream 会持续保持连接，直到返回 false 或客户端断开
	ctx.Stream(func(w io.Writer) bool {
		// 模拟：LLM 逐字生成
		words := []string{"Hello", " ", "I", " ", "am", " ", "Tony", " ", "Bai", "."}

		for _, word := range words {
			// 模拟思考延迟
			time.Sleep(2 * time.Second)

			// 检查客户端是否断开 (Context Done)
			// 这是一个非常重要的细节，防止服务端空转浪费 Token
			select {
			case <-ctx.Done():
				return false // 停止生成
			default:
				// 发送 SSE 事件
				// 格式: data: {"delta": "..."}\n\n
				ctx.SSEvent("message", map[string]string{
					"delta": word,
				})

				// 告诉 Gin 立即将数据推送到网络，不要等待
				ctx.Writer.Flush()
			}
		}

		// 发送结束信号
		ctx.SSEvent("done", "[DONE]")
		return false // 返回 false 结束流
	})
}
