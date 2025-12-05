package errs

import "fmt"

// AppError 实现了 error 接口，承载业务错误信息
type AppError struct {
	Type    ErrorType              // 业务错误码 (机器可读)
	Message string                 // 错误描述 (人类可读)
	Details map[string]interface{} // 结构化详情 (可选，用于扩展)
	Cause   error                  // 原始错误 (用于内部日志，不透传给前端)
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (cause: %v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// 工厂方法：快速创建错误
func New(t ErrorType, msg string) *AppError {
	return &AppError{Type: t, Message: msg}
}

func Wrap(t ErrorType, msg string, cause error) *AppError {
	return &AppError{Type: t, Message: msg, Cause: cause}
}

// WithDetails 添加结构化详情 (支持链式调用)
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}
