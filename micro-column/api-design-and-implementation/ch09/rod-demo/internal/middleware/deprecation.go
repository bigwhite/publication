package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// DeprecationOptions 配置废弃策略
type DeprecationOptions struct {
	SunsetDate time.Time // 彻底下线时间 (HTTP 410 Gone 开始的时间)
	Link       string    // 新版本文档链接或新资源地址
}

// Deprecator 生成一个废弃通知中间件
// 对应 Google AIP-180 及 RFC 8594 标准
func Deprecator(opts DeprecationOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 添加 Deprecation 头
		// 明确告知客户端：此端点已不推荐使用
		c.Header("Deprecation", "true")

		// 2. 添加 Sunset 头 (HTTP Date 格式)
		// 告诉客户端：在这个时间点之后，这个接口将不可用
		if !opts.SunsetDate.IsZero() {
			c.Header("Sunset", opts.SunsetDate.Format(http.TimeFormat))
		}

		// 3. 添加 Link 头，指向新版本
		// rel="successor-version" 是 IANA 注册的标准关联关系，表示"后续版本"
		if opts.Link != "" {
			linkHeader := fmt.Sprintf(`<%s>; rel="successor-version"`, opts.Link)
			c.Header("Link", linkHeader)
		}

		// 继续处理请求
		c.Next()
	}
}
