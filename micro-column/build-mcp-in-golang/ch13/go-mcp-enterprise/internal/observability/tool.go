package observability

import (
	"context"
	"log/slog"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// WrapTool 为 MCP Tool Handler 增加日志和指标采集。
//
// 它不改变原 Handler 的输入输出类型，也不记录具体业务参数，
// 避免把订单号、客户信息等敏感数据写入日志。
func WrapTool[In, Out any](
	logger *slog.Logger,
	metrics *Metrics,
	toolName string,
	next func(
		context.Context,
		*mcp.CallToolRequest,
		In,
	) (*mcp.CallToolResult, Out, error),
) func(
	context.Context,
	*mcp.CallToolRequest,
	In,
) (*mcp.CallToolResult, Out, error) {
	if logger == nil {
		logger = slog.Default()
	}

	return func(
		ctx context.Context,
		req *mcp.CallToolRequest,
		input In,
	) (result *mcp.CallToolResult, output Out, err error) {
		start := time.Now()

		logger.InfoContext(
			ctx,
			"mcp tool started",
			"tool",
			toolName,
		)

		defer func() {
			duration := time.Since(start)

			failed := err != nil
			if result != nil && result.IsError {
				failed = true
			}

			if metrics != nil {
				metrics.Record(toolName, duration, failed)
			}

			outcome := "success"
			if failed {
				outcome = "error"
			}

			attrs := []any{
				"tool", toolName,
				"outcome", outcome,
				"duration_ms", float64(duration.Microseconds()) / 1000,
			}

			if err != nil {
				attrs = append(attrs, "error", err.Error())
				logger.ErrorContext(ctx, "mcp tool completed", attrs...)
				return
			}

			if failed {
				logger.WarnContext(ctx, "mcp tool completed", attrs...)
				return
			}

			logger.InfoContext(ctx, "mcp tool completed", attrs...)
		}()

		result, output, err = next(ctx, req, input)
		return
	}
}
