package middleware

import (
	"fmt"
	"net/http"
	"rod-demo/pkg/limiter"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// KeyStrategy 定义如何提取限流的 Key
type KeyStrategy func(c *gin.Context) string

// 常用策略：基于 IP
func IPKeyStrategy(c *gin.Context) string {
	return "ip:" + c.ClientIP()
}

// 常用策略：基于 API Key (假设在 Header 中)
func APIKeyStrategy(c *gin.Context) string {
	return "apikey:" + c.GetHeader("X-API-Key")
}

// RateLimit 创建限流中间件
// l: 限流器实例
// keyGen: Key 生成策略
// limit: 限流规则
func RateLimit(l *limiter.Limiter, keyGen KeyStrategy, limit limiter.LimitDefinition) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 生成限流 Key
		// 建议加上 request path，实现针对具体接口的限流
		key := fmt.Sprintf("%s:%s", keyGen(c), c.FullPath())

		if c.FullPath() == "" {
			// 404 的情况，也应该限流，防止扫描器
			key = fmt.Sprintf("%s:404", keyGen(c))
		}

		// 2. 检查限流
		res, err := l.Allow(c, key, limit)
		if err != nil {
			// Redis 挂了怎么办？
			// Fail Open: 允许通过 (保障可用性)
			// Fail Closed: 拒绝请求 (保障安全性)
			// 这里演示 Fail Open，记录日志即可
			fmt.Printf("Rate limit redis error: %v\n", err)
			c.Next()
			return
		}

		// 3. 设置标准 Header (无论成功失败都建议设置，让客户端感知配额)
		c.Header("X-RateLimit-Limit", strconv.Itoa(limit.Burst))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(res.Remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(res.ResetAfter.Milliseconds(), 10))

		// 4. 判断结果
		if res.Allowed == 0 {
			// 被限流了
			// 计算 Retry-After (秒)
			retryAfter := int(res.RetryAfter / time.Second)
			if retryAfter < 1 {
				retryAfter = 1
			}
			c.Header("Retry-After", strconv.Itoa(retryAfter))

			// 返回 429
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too Many Requests",
				"retry_after": fmt.Sprintf("%ds", retryAfter),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
