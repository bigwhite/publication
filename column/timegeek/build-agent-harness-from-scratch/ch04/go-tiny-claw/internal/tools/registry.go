// internal/tools/registry.go
package tools

import (
	"context"

	"github.com/yourname/go-tiny-claw/internal/schema"
)

type Registry interface {
	GetAvailableTools() []schema.ToolDefinition
	Execute(ctx context.Context, call schema.ToolCall) schema.ToolResult
}
