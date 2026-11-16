package main

import (
	"log"
	"net/http"

	"example.com/go-mcp-enterprise/internal/authdemo"
	"example.com/go-mcp-enterprise/internal/mcpserver"
	"example.com/go-mcp-enterprise/internal/orders"

	"github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/modelcontextprotocol/go-sdk/oauthex"
)

func main() {
	orderService := orders.NewService(
		orders.NewMemoryRepository(),
	)

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "order-auth-server",
			Version: "v0.1.0",
		},
		nil,
	)

	mcpserver.RegisterAuthorizedOrderTools(server, orderService)

	mcpHandler := mcp.NewStreamableHTTPHandler(
		func(r *http.Request) *mcp.Server {
			return server
		},
		&mcp.StreamableHTTPOptions{
			Stateless: true,
		},
	)

	authMiddleware := auth.RequireBearerToken(
		authdemo.VerifyToken,
		&auth.RequireBearerTokenOptions{
			Scopes: []string{"orders:read"},
		},
	)

	// 本地演示元数据。authorization_servers 使用保留的示例域名，
	// 仅用于展示元数据端点的结构，不代表一个可实际完成授权的服务。
	metadata := &oauthex.ProtectedResourceMetadata{
		Resource: "http://127.0.0.1:8080/mcp",
		AuthorizationServers: []string{
			"https://identity.example.invalid",
		},
		ScopesSupported: []string{
			"orders:read",
			"orders:write",
		},
	}

	http.Handle(
		"/.well-known/oauth-protected-resource",
		auth.ProtectedResourceMetadataHandler(metadata),
	)

	http.Handle("/mcp", authMiddleware(mcpHandler))

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	log.Println("Authenticated MCP Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
