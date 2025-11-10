// pkg/handler/link_handler.go
package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bigwhite/shortlink/pkg/domain"
)

// 定义 service 接口，这是 handler 依赖的契约
type LinkService interface {
	CreateLink(ctx context.Context, originalURL string) (*domain.Link, error)
	Redirect(ctx context.Context, code string) (*domain.Link, error)
	// 假设 GetStats 是 service 层的一个新方法
	GetStats(ctx context.Context, code string) (int64, error)
}

// LinkHandler 持有其依赖
type LinkHandler struct {
	service LinkService
}

func NewLinkHandler(svc LinkService) *LinkHandler {
	return &LinkHandler{service: svc}
}

// CreateLinkRequest 定义了创建链接请求的 JSON 结构
type CreateLinkRequest struct {
	URL string `json:"url"`
}

// CreateLinkResponse 定义了成功创建后的响应 JSON 结构
type CreateLinkResponse struct {
	ShortCode string `json:"short_code"`
}

// CreateLink 是处理创建短链接请求的 http.HandlerFunc
func (h *LinkHandler) CreateLink(w http.ResponseWriter, r *http.Request) {
	var req CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	link, err := h.service.CreateLink(r.Context(), req.URL)
	if err != nil {
		// 在真实项目中，这里应该根据 error 类型返回更精细的状态码
		http.Error(w, "Failed to create link", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateLinkResponse{ShortCode: link.ShortCode})
}

func (h *LinkHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	// 从 URL 路径中提取 short_code
	// 注意：在真实的 Mux 中，我们会用更优雅的方式获取路径参数
	code := strings.TrimPrefix(r.URL.Path, "/")
	if code == "" {
		http.NotFound(w, r)
		return
	}

	link, err := h.service.Redirect(r.Context(), code)
	if err != nil {
		// 在真实项目中，这里应该根据 error 类型返回更精细的状态码
		// 比如，如果是 "link not found"，应该返回 404
		http.Error(w, "Link not found or internal error", http.StatusNotFound)
		return
	}

	// 执行重定向
	http.Redirect(w, r, link.OriginalURL, http.StatusFound) // 302 Found
}

// GetStatsResponse 定义了统计接口的响应
type GetStatsResponse struct {
	Visits int64 `json:"visits"`
}

func (h *LinkHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	// 同样，从路径中提取 short_code
	// e.g., /api/links/{short_code}/stats
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 5 {
		http.NotFound(w, r)
		return
	}
	code := parts[3]

	visits, err := h.service.GetStats(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to get stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(GetStatsResponse{Visits: visits})
}


