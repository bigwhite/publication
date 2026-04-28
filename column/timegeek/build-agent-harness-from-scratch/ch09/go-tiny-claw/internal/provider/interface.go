// internal/provider/interface.go
package provider

import (
	"context"

	"github.com/yourname/go-tiny-claw/internal/schema"
)

// LLMProvider defines the unified interface for communicating with large models
type LLMProvider interface {
	// Generate receives the current context history and available tools list, returns the model response
	Generate(ctx context.Context, messages []schema.Message, availableTools []schema.ToolDefinition) (*schema.Message, error)
}
