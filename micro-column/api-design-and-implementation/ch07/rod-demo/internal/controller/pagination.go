package controller

// ListRequest 定义符合 AIP-158 的标准分页请求参数
type ListRequest struct {
	PageSize  int    `form:"page_size"`  // 对应 ?page_size=10
	PageToken string `form:"page_token"` // 对应 ?page_token=abc...
}

// ListResponse 定义符合 AIP-158 的标准分页响应结构
// T 是资源类型
type ListResponse[T any] struct {
	Items         []T    `json:"items"`                // 资源列表
	NextPageToken string `json:"next_page_token"`      // 下一页游标，为空表示结束
	TotalSize     int    `json:"total_size,omitempty"` // 可选：总条数
}
