package middleware

import (
	"bytes"

	"context"
	"encoding/json"
	"net/http"
	"rod-demo/pkg/redis"
	"time"

	"github.com/gin-gonic/gin"
)

type bodyDumpResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	// 1. 写入缓存
	w.body.Write(b)
	// 2. 真正写入网络响应
	return w.ResponseWriter.Write(b)
}

func (w *bodyDumpResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

const (
	HeaderIdempotencyKey = "X-Idempotency-Key"
	KeyPrefix            = "idemp:"
	LockExpire           = 24 * time.Hour // 幂等键有效期，通常设为24小时
)

// 缓存的响应结构
type cachedResponse struct {
	Status  int                 `json:"status"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

func Idempotency() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 0. 方法过滤
		// 根据 HTTP 语义，GET, HEAD, OPTIONS, DELETE 天然幂等（或者只读），无需幂等键保护
		// 这里的策略取决于你的业务需求：
		// - GET: 绝对跳过
		// - DELETE: 标准语义是幂等的，但如果业务逻辑重，也可以加锁。这里遵循注释建议跳过。
		if c.Request.Method == http.MethodGet ||
			c.Request.Method == http.MethodHead ||
			c.Request.Method == http.MethodOptions ||
			c.Request.Method == http.MethodDelete {

			c.Next()
			return
		}

		// 1. 获取幂等键
		key := c.GetHeader(HeaderIdempotencyKey)
		if key == "" {
			// 如果没有传 key，则不进行幂等保护，直接透传
			c.Next()
			return
		}

		redisKey := KeyPrefix + key
		ctx := context.Background()

		// 2. 检查 Key 是否存在 (使用 SETNX 实现原子锁)
		// 状态定义：
		// - 不存在: 第一次请求，抢锁成功
		// - 存在: 重复请求

		// 这里我们尝试抢锁，为了简单，我们用 Redis String 存储状态
		// 实际生产可以使用 Lua 脚本保证更复杂的原子性

		// 尝试获取结果缓存
		val, err := redis.Client.Get(ctx, redisKey).Result()

		if err == nil {
			// 2.1 Case A: Key 存在
			if val == "PROCESSING" {
				// 正在处理中，发生并发请求
				c.JSON(http.StatusConflict, gin.H{
					"error": "Duplicate request, processing in progress",
				})
				c.Abort()
				return
			}

			// 处理完毕，直接返回缓存的结果
			var resp cachedResponse
			json.Unmarshal([]byte(val), &resp)

			// 恢复 Header
			for k, v := range resp.Headers {
				for _, s := range v {
					c.Header(k, s)
				}
			}
			// 恢复 Status 和 Body
			c.Data(resp.Status, c.Writer.Header().Get("Content-Type"), []byte(resp.Body))
			c.Abort() // 阻断后续 Handler 执行
			return
		}

		// 2.2 Case B: Key 不存在，我们需要抢锁
		// SETNX: Set if Not Exists
		success, _ := redis.Client.SetNX(ctx, redisKey, "PROCESSING", LockExpire).Result()
		if !success {
			// 抢锁失败（极低概率并发），视为冲突
			c.JSON(http.StatusConflict, gin.H{"error": "Concurrent request conflict"})
			c.Abort()
			return
		}

		// 3. 包装 ResponseWriter 以捕获响应
		dumpWriter := &bodyDumpResponseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = dumpWriter

		// 4. 执行业务逻辑
		c.Next()

		// 5. 业务执行完毕，缓存响应结果
		// 注意：只缓存成功的或特定的错误码，视业务需求而定
		if c.Writer.Status() < 500 {
			cacheObj := cachedResponse{
				Status:  c.Writer.Status(),
				Headers: c.Writer.Header(),
				Body:    dumpWriter.body.String(),
			}
			jsonBytes, _ := json.Marshal(cacheObj)

			// 更新 Redis，将 "PROCESSING" 替换为 真实结果
			redis.Client.Set(ctx, redisKey, string(jsonBytes), LockExpire)
		} else {
			// 如果业务报错（500），通常应该删除 Key，允许客户端重试
			redis.Client.Del(ctx, redisKey)
		}
	}
}
