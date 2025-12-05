package middleware

import (
	"rod-demo/pkg/errs"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 执行后续中间件和业务逻辑
		c.Next()

		// 2. 检查是否有错误产生
		// Gin 允许在 Handler 中通过 c.Error(err) 收集错误
		if len(c.Errors) == 0 {
			return
		}

		// 3.以此请求中最后一个错误为准
		lastErr := c.Errors.Last().Err

		var appErr *errs.AppError
		// 4. 类型断言：如果是我们定义的 AppError
		if e, ok := lastErr.(*errs.AppError); ok {
			appErr = e
		} else {
			// 5. 如果是未知错误 (如 panic 或第三方库错误)，包装为 Internal Server Error
			// 实际生产中这里应该打印堆栈日志
			appErr = errs.Wrap(errs.ErrInternalServer, "something went wrong", lastErr)
		}

		// 6. 构造 RFC 7807 响应
		problem := errs.ProblemDetails{
			Type:     "https://example.com/probs/" + string(appErr.Type), // 示例 URI
			Title:    string(appErr.Type),
			Status:   appErr.Type.HTTPStatus(),
			Detail:   appErr.Message,
			Instance: c.Request.RequestURI,
			Details:  appErr.Details,
		}

		// 7. 发送响应
		// 注意：使用 application/problem+json 作为 Content-Type 是 RFC 推荐的
		c.Header("Content-Type", "application/problem+json")
		c.JSON(problem.Status, problem)
	}
}
