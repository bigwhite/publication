package domain

// Operation 遵循 Google AIP-151 标准结构
type Operation struct {
	ID       string      `json:"id"`       // 对应 name: operations/{id}
	Done     bool        `json:"done"`     // 任务是否完成
	Metadata interface{} `json:"metadata"` // 进度或上下文
	Response interface{} `json:"response"` // 成功返回值
	Error    interface{} `json:"error"`    // 失败返回值
}

// OperationMetadata 具体的进度信息
type OperationMetadata struct {
	Progress int    `json:"progress_percent"`
	Status   string `json:"status"`
}
