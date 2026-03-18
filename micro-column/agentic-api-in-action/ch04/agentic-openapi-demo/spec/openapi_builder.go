package spec

import (
	"sync"
)

// AgenticMeta 定义了我们注入到 OpenAPI 中的智能体元数据
type AgenticMeta struct {
	Action        string   `json:"x-agent-action"`
	Preconditions []string `json:"x-agent-preconditions,omitempty"`
	SideEffects   []string `json:"x-agent-side-effects,omitempty"`
	Hints         string   `json:"x-agent-hints,omitempty"`
}

// Operation 定义了 OpenAPI 中的单个操作 (如 GET /users)
type Operation struct {
	Summary     string                 `json:"summary"`
	OperationID string                 `json:"operationId"`
	Parameters  []interface{}          `json:"parameters,omitempty"`
	RequestBody interface{}            `json:"requestBody,omitempty"`
	Responses   map[string]interface{} `json:"responses"`
	// 嵌入 Agentic 元数据，注意 JSON 序列化时它会被平铺到当前层级
	AgenticMeta
}

// OpenAPISpec 维护全局的 API 规范
type OpenAPISpec struct {
	mu    sync.RWMutex
	Paths map[string]map[string]Operation `json:"paths"`
}

var GlobalSpec = &OpenAPISpec{
	Paths: make(map[string]map[string]Operation),
}

// RegisterRoute 动态注册一个路由及其 OpenAPI 定义
func (s *OpenAPISpec) RegisterRoute(path string, method string, op Operation) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Paths[path] == nil {
		s.Paths[path] = make(map[string]Operation)
	}
	s.Paths[path][method] = op
}

// GenerateJSON 生成最终的 OpenAPI 文档 (简化版)
func (s *OpenAPISpec) GenerateJSON() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	doc := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":       "Agentic API System",
			"version":     "1.0.0",
			"description": "This API is optimized for AI Agents with semantic meta-data.",
		},
		"paths": s.Paths,
	}
	return doc
}
