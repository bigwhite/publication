package registry

import (
	"sync"
)

// ActionCategory 定义了我们上一讲学过的六大分类
type ActionCategory string

const (
	CategoryAcquire     ActionCategory = "ACQUIRE"
	CategoryCompute     ActionCategory = "COMPUTE"
	CategoryTransact    ActionCategory = "TRANSACT"
	CategoryIntegrate   ActionCategory = "INTEGRATE"
	CategoryOrchestrate ActionCategory = "ORCHESTRATE"
	CategoryNotify      ActionCategory = "NOTIFY"
)

// Capability 代表系统向 AI 暴露的一个具体能力
type Capability struct {
	Name          string         `json:"name"`           // 动作名称，如 "refund_order"
	Category      ActionCategory `json:"category"`       // 所属分类，如 "TRANSACT"
	Description   string         `json:"description"`    // 给 AI 看的详细描述
	Endpoint      string         `json:"endpoint"`       // 实际调用的 URL 路径
	Preconditions []string       `json:"preconditions"`  // 必须满足的前置条件 (自然语言描述)
	RequiredScope string         `json:"required_scope"` // 执行所需的权限标识
}

// CapabilityRegistry 是全局的能力注册中心
type CapabilityRegistry struct {
	mu           sync.RWMutex
	capabilities map[string]Capability
}

// NewRegistry 创建一个新的注册表
func NewRegistry() *CapabilityRegistry {
	return &CapabilityRegistry{
		capabilities: make(map[string]Capability),
	}
}

// Register 注册一个新的能力
func (r *CapabilityRegistry) Register(cap Capability) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.capabilities[cap.Name] = cap
}

// ListAll 返回所有已注册的能力，供 DISCOVER 接口使用
func (r *CapabilityRegistry) ListAll() []Capability {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var list []Capability
	for _, cap := range r.capabilities {
		list = append(list, cap)
	}
	return list
}

// GlobalRegistry 用于本示例的简化访问
var Global = NewRegistry()
