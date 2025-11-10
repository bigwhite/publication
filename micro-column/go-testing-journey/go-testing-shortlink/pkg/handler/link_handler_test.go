//go:build unit

package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bigwhite/shortlink/pkg/domain"
	"github.com/bigwhite/shortlink/pkg/handler"
)

// 1. 创建 Stub Object
// StubLinkService 是 LinkService 接口的一个测试替身
type StubLinkService struct {
	CreateLinkFunc func(ctx context.Context, originalURL string) (*domain.Link, error)
	RedirectFunc   func(ctx context.Context, code string) (*domain.Link, error)
	GetStatsFunc   func(ctx context.Context, code string) (int64, error)
}

func (s *StubLinkService) Redirect(ctx context.Context, code string) (*domain.Link, error) {
	if s.RedirectFunc != nil {
		return s.RedirectFunc(ctx, code)
	}
	return nil, errors.New("RedirectFunc not implemented")
}

func (s *StubLinkService) CreateLink(ctx context.Context, originalURL string) (*domain.Link, error) {
	if s.CreateLinkFunc != nil {
		return s.CreateLinkFunc(ctx, originalURL)
	}
	return nil, errors.New("CreateLinkFunc not implemented")
}

func (s *StubLinkService) GetStats(ctx context.Context, code string) (int64, error) {
	if s.GetStatsFunc != nil {
		return s.GetStatsFunc(ctx, code)
	}
	return 0, errors.New("GetStatsFunc not implemented")
}

func TestLinkHandler_CreateLink(t *testing.T) {
	// 2. 使用表驱动测试
	testCases := []struct {
		name           string
		reqBody        string
		stub           *StubLinkService
		wantStatusCode int
		wantRespBody   string
	}{
		{
			name:    "成功创建",
			reqBody: `{"url": "https://example.com"}`,
			stub: &StubLinkService{
				CreateLinkFunc: func(ctx context.Context, originalURL string) (*domain.Link, error) {
					return &domain.Link{ShortCode: "success"}, nil
				},
			},
			wantStatusCode: http.StatusCreated,
			wantRespBody:   `{"short_code":"success"}`,
		},
		{
			name:           "请求体 JSON 格式错误",
			reqBody:        `{"url": "https://example.com"`, // 缺少右括号
			stub:           &StubLinkService{},              // 不会被调用
			wantStatusCode: http.StatusBadRequest,
			wantRespBody:   `Invalid request body`,
		},
		{
			name:    "Service 返回错误",
			reqBody: `{"url": "https://example.com"}`,
			stub: &StubLinkService{
				CreateLinkFunc: func(ctx context.Context, originalURL string) (*domain.Link, error) {
					return nil, errors.New("internal error")
				},
			},
			wantStatusCode: http.StatusInternalServerError,
			wantRespBody:   `Failed to create link`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 3. 准备请求和响应记录器
			req := httptest.NewRequest("POST", "/api/links", strings.NewReader(tc.reqBody))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			// 4. 注入 Stub 并创建 Handler
			linkHandler := handler.NewLinkHandler(tc.stub)

			// 5. 执行 Handler
			http.HandlerFunc(linkHandler.CreateLink).ServeHTTP(rr, req)

			// 6. 断言结果
			if rr.Code != tc.wantStatusCode {
				t.Errorf("status code got %d, want %d", rr.Code, tc.wantStatusCode)
			}

			// 对 Body 的断言需要注意，错误响应可能包含换行符
			trimmedBody := strings.TrimSpace(rr.Body.String())

			// 如果期望是 JSON，我们可以反序列化后比较，更健壮
			if strings.HasPrefix(tc.wantRespBody, "{") {
				var got, want map[string]interface{}
				if err := json.Unmarshal([]byte(trimmedBody), &got); err != nil {
					t.Fatalf("failed to unmarshal response body: %v", err)
				}
				if err := json.Unmarshal([]byte(tc.wantRespBody), &want); err != nil {
					t.Fatalf("failed to unmarshal wantRespBody: %v", err)
				}
				// 这里可以用更完善的 deep equal 库
				if got["short_code"] != want["short_code"] {
					t.Errorf("response body got %v, want %v", got, want)
				}
			} else {
				if trimmedBody != tc.wantRespBody {
					t.Errorf("response body got %q, want %q", trimmedBody, tc.wantRespBody)
				}
			}
		})
	}
}
