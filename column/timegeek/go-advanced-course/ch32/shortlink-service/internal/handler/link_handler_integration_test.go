// internal/handler/link_handler_integration_test.go
package handler_test // 使用 _test 包名，表示从包外部进行测试 (黑盒)

import (
	"bytes" // 引入context
	"encoding/json"
	"io" // 用于 slog 的 io.Discard
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	// 用于 context.WithTimeout
	// 导入被测试的包和依赖 (路径根据你的go.mod)
	"github.com/your_org/shortlink/internal/handler"      // 我们的handler
	"github.com/your_org/shortlink/internal/service"      // 我们的service
	"github.com/your_org/shortlink/internal/store/memory" // 使用内存存储

	// "github.com/your_org/shortlink/internal/store" // 如果需要 store.ErrNotFound

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel" // 为了service能获取tracer
)

// setupIntegrationTestServer 创建一个包含真实（内存）依赖的测试服务器。
// 返回一个 http.Handler，可以直接用于 httptest.Server 或直接调用其 ServeHTTP。
func setupIntegrationTestServer(t *testing.T) http.Handler {
	// 1. 创建 Logger (在测试中，我们可能不关心日志输出，使用Discard Handler)
	testLogger := slog.New(slog.NewTextHandler(io.Discard, nil)) // 忽略所有日志输出

	// 2. 创建真实的内存 Store 实例
	memStore := memory.NewStore(testLogger.With("component", "memory_store_integ_test"))
	// t.Cleanup(func() { memStore.Close() }) // 如果 memStore 需要清理

	// 3. 创建真实的 Service 实例，注入内存 Store 和 Logger
	// 假设NewShortenerService现在也接收一个tracer，但对于集成测试，我们可以用NoopTracer
	// otel.SetTracerProvider(trace.NewNoopTracerProvider()) // 确保有一个Provider
	// tracer := otel.Tracer("integ-test-tracer")
	// shortenerSvc := service.NewShortenerService(memStore, testLogger.With("component", "service_integ_test"), tracer)
	// For simplicity, let's assume our existing NewShortenerService in service.go can be used.
	// If it strictly requires a global tracer to be set, TestMain or a test setup func should do it.
	// For this test, we'll re-use the NewShortenerService from service package.
	// Make sure the NewShortenerService in `service` package is compatible.
	// It expects (store.Store, *slog.Logger). The tracer is acquired via otel.Tracer() internally.

	// Ensure a global tracer provider is set for the service to pick up a tracer,
	// even if it's a NoOp one for tests not focusing on tracing.
	// (This would ideally be in a TestMain or a setup helper for all integration tests)
	// For now, we ensure it's callable.
	if otel.GetTracerProvider() == nil {
		// In a real setup, you might initialize a NoOp tracer provider here for tests
		// or ensure your InitTracerProvider from tracing package is test-friendly.
		// For this example, we'll assume the service's otel.Tracer() call will get a NoOp if none is set.
	}

	shortenerSvc := service.NewShortenerService(memStore, testLogger.With("component", "service_integ_test"), nil)

	// 4. 创建真实的 Handler 实例，注入 Service 和 Logger
	linkHdlr := handler.NewLinkHandler(shortenerSvc, testLogger.With("component", "handler_integ_test"))

	// 5. 创建一个 HTTP Mux 并注册路由 (与app.go中类似，但更简化)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/links", linkHdlr.CreateShortLink) // 路径与app.go中一致
	// 模拟 /{shortCode} 的重定向路由
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && len(r.URL.Path) > 1 && !strings.HasPrefix(r.URL.Path, "/api/") {
			shortCode := strings.TrimPrefix(r.URL.Path, "/")
			linkHdlr.RedirectShortLink(w, r, shortCode) // 调用 RedirectShortLink
		} else {
			http.NotFound(w, r)
		}
	})
	return mux
}

func TestIntegration_CreateAndRedirectLink(t *testing.T) {
	// 1. 设置测试服务器
	testServerHandler := setupIntegrationTestServer(t)

	// 定义测试场景的输入
	longURL := "https://www.example.com/a-very-long-url-for-integration-testing"
	createPayload := handler.CreateLinkRequest{LongURL: longURL} // 使用handler中定义的结构体
	payloadBytes, err := json.Marshal(createPayload)
	require.NoError(t, err, "Failed to marshal create link request payload")

	// --- 步骤1: 创建短链接 ---
	reqCreate := httptest.NewRequest(http.MethodPost, "/api/links", bytes.NewReader(payloadBytes))
	reqCreate.Header.Set("Content-Type", "application/json")
	rrCreate := httptest.NewRecorder() // httptest.ResponseRecorder 用于捕获响应

	testServerHandler.ServeHTTP(rrCreate, reqCreate) //直接调用Handler的ServeHTTP

	// 断言创建请求的响应
	require.Equal(t, http.StatusCreated, rrCreate.Code, "CreateLink: Unexpected status code")

	var createResp handler.CreateLinkResponse // 使用handler中定义的结构体
	err = json.Unmarshal(rrCreate.Body.Bytes(), &createResp)
	require.NoError(t, err, "CreateLink: Failed to unmarshal response body")
	require.NotEmpty(t, createResp.ShortCode, "CreateLink: Short code in response should not be empty")

	shortCodeGenerated := createResp.ShortCode
	t.Logf("CreateLink: Successfully created short code '%s' for URL '%s'", shortCodeGenerated, longURL)

	// --- 步骤2: 使用生成的短链接进行重定向 ---
	// 这里我们直接访问 /{shortCode} 路径
	redirectPath := "/" + shortCodeGenerated
	reqRedirect := httptest.NewRequest(http.MethodGet, redirectPath, nil)
	// 为了能正确获取 Location header，我们需要一个能处理重定向的客户端，
	// 或者检查 ResponseRecorder 的 Header。httptest.Recorder 不会自动跟随重定向。
	rrRedirect := httptest.NewRecorder()

	testServerHandler.ServeHTTP(rrRedirect, reqRedirect)

	// 断言重定向请求的响应
	// 对于短链接服务，我们通常期望301 (永久)或302 (临时/找到)重定向
	// 在我们的 RedirectShortLink handler 中，我们使用了 http.StatusFound (302)
	assert.Equal(t, http.StatusFound, rrRedirect.Code, "RedirectShortLink: Unexpected status code")

	// 检查 Location 响应头是否指向原始的长链接
	redirectLocation := rrRedirect.Header().Get("Location")
	assert.Equal(t, longURL, redirectLocation, "RedirectShortLink: Redirect Location mismatch")
	t.Logf("RedirectShortLink: Successfully redirected from '%s' to '%s'", redirectPath, redirectLocation)

	// --- (可选) 步骤3: 尝试获取一个不存在的短链接 ---
	reqNotFound := httptest.NewRequest(http.MethodGet, "/nonexistentcode", nil)
	rrNotFound := httptest.NewRecorder()
	testServerHandler.ServeHTTP(rrNotFound, reqNotFound)
	assert.Equal(t, http.StatusNotFound, rrNotFound.Code, "GetNonExistentLink: Expected HTTP 404 Not Found")
	t.Logf("GetNonExistentLink: Correctly received 404 for '/nonexistentcode'")
}

// 可以添加更多集成测试用例，例如：
// - 测试无效输入 (如创建时long_url为空)
// - 测试并发创建 (如果store支持并发)
// - 测试短码大小写不敏感 (如果业务逻辑如此定义)
// - 等等

func TestIntegration_CreateLink_InvalidInput(t *testing.T) {
	testServerHandler := setupIntegrationTestServer(t)

	t.Run("EmptyLongURL", func(t *testing.T) {
		createPayload := handler.CreateLinkRequest{LongURL: ""}
		payloadBytes, _ := json.Marshal(createPayload)
		reqCreate := httptest.NewRequest(http.MethodPost, "/api/links", bytes.NewReader(payloadBytes))
		reqCreate.Header.Set("Content-Type", "application/json")
		rrCreate := httptest.NewRecorder()

		testServerHandler.ServeHTTP(rrCreate, reqCreate)
		assert.Equal(t, http.StatusBadRequest, rrCreate.Code, "Expected 400 Bad Request for empty long_url")
	})
	// 可以添加其他无效输入的测试用例
}
