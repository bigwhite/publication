package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/your_org/shortlink/internal/shortener" // 导入 shortener 包以使用其错误和类型
)

// LinkAPI 封装了与短链接相关的HTTP处理逻辑
type LinkAPI struct {
	service *shortener.Service
	logger  *log.Logger
}

// NewLinkAPI 创建一个新的 LinkAPI 处理器
func NewLinkAPI(svc *shortener.Service, baseLogger *log.Logger) *LinkAPI {
	if baseLogger == nil {
		baseLogger = log.New(os.Stdout, "[LinkAPI-Default] ", log.LstdFlags|log.Lshortfile)
	}
	handlerLogger := log.New(baseLogger.Writer(), fmt.Sprintf("%s[LinkAPIHandler] ", baseLogger.Prefix()), baseLogger.Flags())
	return &LinkAPI{
		service: svc,
		logger:  handlerLogger,
	}
}

// CreateShortLinkRequest 是创建短链接请求的结构体
type CreateShortLinkRequest struct {
	LongURL string `json:"long_url"`
}

// CreateShortLinkResponse 是创建短链接响应的结构体
type CreateShortLinkResponse struct {
	ShortCode string `json:"short_code"`
}

// preview 是一个辅助函数，用于截断长字符串以便日志记录
func preview(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}

// CreateLink 处理创建短链接的HTTP POST请求
func (h *LinkAPI) CreateLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		h.logger.Printf("WARN: Invalid method for create link: %s from %s\n", r.Method, r.RemoteAddr)
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateShortLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Printf("ERROR: Failed to decode request body from %s: %v\n", r.RemoteAddr, err)
		http.Error(w, `{"error": "Invalid request body: `+err.Error()+`"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if strings.TrimSpace(req.LongURL) == "" {
		h.logger.Printf("WARN: Long URL is empty in request from %s\n", r.RemoteAddr)
		http.Error(w, `{"error": "long_url cannot be empty"}`, http.StatusBadRequest)
		return
	}

	h.logger.Printf("INFO: Received request to create short link from %s, LongURL_preview: %s\n", r.RemoteAddr, preview(req.LongURL, 50)) // 使用本地 preview

	shortCode, err := h.service.CreateShortLink(ctx, req.LongURL)
	if err != nil {
		h.logger.Printf("ERROR: Service failed to create short link for %s from %s: %v\n", preview(req.LongURL, 50), r.RemoteAddr, err) // 使用本地 preview
		if errors.Is(err, shortener.ErrInvalidLongURL) {
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		} else if errors.Is(err, shortener.ErrConflict) || errors.Is(err, shortener.ErrShortCodeGenerationFailed) { // 使用 shortener.ErrConflict
			http.Error(w, fmt.Sprintf(`{"error": "Could not create short link: %s"}`, err.Error()), http.StatusConflict)
		} else {
			http.Error(w, `{"error": "Failed to create short link, internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	resp := CreateShortLinkResponse{ShortCode: shortCode}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Printf("ERROR: Failed to encode response for %s: %v\n", r.RemoteAddr, err)
	}
	h.logger.Printf("INFO: Successfully created short link for %s. ShortCode: %s\n", r.RemoteAddr, shortCode)
}

// RedirectLink 处理短链接重定向的HTTP GET请求
func (h *LinkAPI) RedirectLink(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shortCode := strings.TrimPrefix(r.URL.Path, "/")
	h.logger.Printf("INFO: Received request to redirect short link from %s. ShortCode: %s, Path: %s\n", r.RemoteAddr, shortCode, r.URL.Path)

	if r.Method != http.MethodGet {
		h.logger.Printf("WARN: Invalid method for redirect link: %s from %s\n", r.Method, r.RemoteAddr)
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	if strings.TrimSpace(shortCode) == "" || shortCode == "api/links" || shortCode == "healthz" {
		h.logger.Printf("INFO: Path is not a shortcode, treating as not found. Path: %s, from %s\n", r.URL.Path, r.RemoteAddr)
		http.NotFound(w, r)
		return
	}

	longURL, err := h.service.GetAndTrackLongURL(ctx, shortCode)
	if err != nil {
		h.logger.Printf("WARN: Service failed to get long URL for redirect from %s. ShortCode: %s, Error: %v\n", r.RemoteAddr, shortCode, err)
		if errors.Is(err, shortener.ErrLinkNotFound) || errors.Is(err, shortener.ErrShortCodeTooShort) { // 使用 shortener 定义的错误
			http.Error(w, "Short link not found or invalid", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	h.logger.Printf("INFO: Redirecting %s from %s to %s\n", shortCode, r.RemoteAddr, preview(longURL, 50)) // 使用本地 preview
	http.Redirect(w, r, longURL, http.StatusFound)
}
