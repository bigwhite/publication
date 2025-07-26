// pkg/lifecycle/lifecycle.go
package lifecycle

import "context"

// Component 定义了可管理的应用组件的接口，
// 这些组件具有明确的启动和停止生命周期。
type Component interface {
	Start(ctx context.Context) error // 启动组件。
	Stop(ctx context.Context) error  // 优雅地停止组件。
	Name() string                    // 返回组件名称，用于日志记录。
}
