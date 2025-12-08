package errs

// ProblemDetails 符合 RFC 7807 的 JSON 结构
type ProblemDetails struct {
	Type     string                 `json:"type"`               // 错误类型的 URI 标识
	Title    string                 `json:"title"`              // 简短描述 (对应 ErrorType)
	Status   int                    `json:"status"`             // HTTP 状态码
	Detail   string                 `json:"detail"`             // 详细描述 (对应 Message)
	Instance string                 `json:"instance,omitempty"` // 请求路径
	Details  map[string]interface{} `json:"details,omitempty"`  // 扩展字段 (Google AIP 风格)
}
