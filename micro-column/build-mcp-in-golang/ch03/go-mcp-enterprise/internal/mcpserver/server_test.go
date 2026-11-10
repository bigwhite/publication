package mcpserver

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestServerConnectAndListTools(t *testing.T) {
	ctx := context.Background()

	logger := slog.New(
		slog.NewTextHandler(io.Discard, nil),
	)

	server := New(logger)

	// 创建一对内存 Transport。
	serverTransport, clientTransport := mcp.NewInMemoryTransports()

	// 先连接 Server。
	serverSession, err := server.Connect(
		ctx,
		serverTransport,
		nil,
	)
	if err != nil {
		t.Fatalf("connect server: %v", err)
	}

	client := mcp.NewClient(
		&mcp.Implementation{
			Name:    "test-client",
			Version: "0.1.0",
		},
		nil,
	)

	// 再连接 Client。
	clientSession, err := client.Connect(
		ctx,
		clientTransport,
		nil,
	)
	if err != nil {
		t.Fatalf("connect client: %v", err)
	}

	// 关闭 Client 后，等待 Server 侧连接结束。
	t.Cleanup(func() {
		if err := clientSession.Close(); err != nil {
			t.Errorf("close client session: %v", err)
		}

		if err := serverSession.Wait(); err != nil {
			t.Errorf("wait server session: %v", err)
		}
	})

	// 查询 Server 当前提供的 Tool。
	result, err := clientSession.ListTools(
		ctx,
		&mcp.ListToolsParams{},
	)
	if err != nil {
		t.Fatalf("list tools: %v", err)
	}

	if len(result.Tools) != 0 {
		t.Fatalf(
			"expected no tools, got %d",
			len(result.Tools),
		)
	}
}
