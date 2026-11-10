package mcpserver

import (
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// New 创建企业订单系统的 MCP Server。
func New(logger *slog.Logger) *mcp.Server {
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "go-mcp-enterprise",
			Version: "0.1.0",
		},
		&mcp.ServerOptions{
			Instructions: "这是一个企业订单系统 MCP Server。",
			Logger:       logger,
			// 明确设置能力集合，避免继承历史默认能力。
			Capabilities: &mcp.ServerCapabilities{},
		},
	)

	return server
}
