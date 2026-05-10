package engine

import "context"

type Reporter interface {
	OnThinking(ctx context.Context)
	OnToolCall(ctx context.Context, toolName string, args string)
	OnToolResult(ctx context.Context, toolName string, result string, isError bool)
	OnMessage(ctx context.Context, content string)
}
