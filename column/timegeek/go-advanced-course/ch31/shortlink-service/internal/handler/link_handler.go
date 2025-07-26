package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/your_org/shortlink/internal/service"
	"github.com/your_org/shortlink/internal/store" // 需要 ErrNotFound

	"go.opentelemetry.io/otel/attribute" // For Tracing attributes
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace" // For Tracing
)

// CreateLinkRequest 定义了创建短链接请求的JSON结构体
type CreateLinkRequest struct {
	LongURL     string `json:"long_url"`
	UserID      string `json:"user_id,omitempty"`      // (可选)
	OriginalURL string `json:"original_url,omitempty"` // (可选)
	// ExpireInDays int    `json:"expire_in_days,omitempty"` // (可选)
}

// CreateLinkResponse 定义了创建短链接响应的JSON结构体
type CreateLinkResponse struct {
	ShortCode string `json:"short_code"`
	LongURL   string `json:"long_url"` // (可选) 返回原始长链接，方便客户端确认
}

// LinkHandler 负责处理与短链接相关的HTTP请求
type LinkHandler struct {
	svc    service.ShortenerService // 业务逻辑服务接口
	logger *slog.Logger             // 注入的logger实例
	// tracer trace.Tracer // Handler 通常不直接创建顶层Span，而是由中间件处理
	// 但如果Handler内部有多个逻辑块需要细分追踪，可以获取tracer创建子Span
}

// NewLinkHandler 创建一个新的LinkHandler
func NewLinkHandler(svc service.ShortenerService, logger *slog.Logger) *LinkHandler {
	return &LinkHandler{
		svc:    svc,
		logger: logger.With("component", "LinkHandler"), // 为这个组件的logger添加固定属性
		// tracer: otel.Tracer("github.com/your_org/shortlink/internal/handler"), // 如果需要手动创建span
	}
}

// CreateShortLink 处理创建新短链接的请求 (POST /api/links)
func (h *LinkHandler) CreateShortLink(w http.ResponseWriter, r *http.Request) {
	// 从请求的context中获取由OTel HTTP中间件创建的Span
	// 这样后续的日志和手动创建的子Span都能关联到这个请求的Trace
	ctx := r.Context()
	span := trace.SpanFromContext(ctx) // 获取当前Span

	// 为每个请求创建一个上下文相关的logger，加入请求特定的属性
	// (TraceID等通常由日志中间件或slog的Handler自动从ctx中获取并添加)
	requestLogger := h.logger.With(
		slog.String("http_method", r.Method),
		slog.String("http_path", r.URL.Path),
		// 示例：从context中获取trace ID并添加到日志 (需要slog的handler支持或手动添加)
		// slog.String("trace_id", span.SpanContext().TraceID().String()),
	)
	requestLogger.DebugContext(ctx, "Handler: Received request to create short link.")
	span.AddEvent("Handler: Received request") // 在Span上记录事件

	if r.Method != http.MethodPost {
		requestLogger.WarnContext(ctx, "Handler: Invalid HTTP method for CreateShortLink.", slog.String("received_method", r.Method))
		span.SetStatus(codes.Error, "Invalid HTTP method")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqPayload CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&reqPayload); err != nil {
		requestLogger.ErrorContext(ctx, "Handler: Failed to decode request body.", slog.Any("error", err))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to decode request body")
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if strings.TrimSpace(reqPayload.LongURL) == "" {
		requestLogger.WarnContext(ctx, "Handler: LongURL is empty in request.")
		span.SetStatus(codes.Error, "LongURL is empty")
		http.Error(w, "long_url cannot be empty", http.StatusBadRequest)
		return
	}
	span.SetAttributes(attribute.String("request.long_url", reqPayload.LongURL))
	if reqPayload.UserID != "" {
		span.SetAttributes(attribute.String("request.user_id", reqPayload.UserID))
	}

	requestLogger.InfoContext(ctx, "Handler: Processing create short link request.", slog.String("long_url", reqPayload.LongURL))

	// 为service调用设置一个独立的超时（可选，也可以由Service层自己管理）
	// serviceCtx, serviceCancel := context.WithTimeout(ctx, 5*time.Second)
	// defer serviceCancel()
	// 使用请求的ctx，让OTel的span传播
	serviceCtx := ctx

	// 默认过期时间 (可以从配置读取或作为请求参数)
	expireAt := time.Now().Add(time.Hour * 24 * 30) // 默认30天

	shortCode, err := h.svc.CreateShortLink(serviceCtx, reqPayload.LongURL, reqPayload.UserID, reqPayload.OriginalURL, expireAt)
	if err != nil {
		requestLogger.ErrorContext(ctx, "Handler: Service failed to create short link.",
			slog.Any("error", err),
			slog.String("long_url", reqPayload.LongURL))
		span.RecordError(err)
		span.SetStatus(codes.Error, "Service failed to create short link")

		// 根据错误类型返回不同的HTTP状态码
		if errors.Is(err, service.ErrIDGenerationFailed) {
			http.Error(w, "Failed to generate unique short code, please try again.", http.StatusConflict)
		} else if errors.Is(err, errors.New("long URL cannot be empty")) { // 假设service也会校验
			http.Error(w, "long_url cannot be empty", http.StatusBadRequest)
		} else {
			http.Error(w, "Internal server error creating short link.", http.StatusInternalServerError)
		}
		return
	}

	response := CreateLinkResponse{
		ShortCode: shortCode,
		LongURL:   reqPayload.LongURL, // 返回原始长链接，方便客户端
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// 如果Encode失败，此时WriteHeader已发送，只能记录错误了
		requestLogger.ErrorContext(ctx, "Handler: Failed to encode response body.", slog.Any("error", err))
		span.RecordError(err) // 也可以在主span上记录这个错误
	}

	requestLogger.InfoContext(ctx, "Handler: Successfully created and responded with short link.",
		slog.String("short_code", shortCode),
		slog.String("long_url", reqPayload.LongURL),
	)
	span.SetAttributes(attribute.String("response.short_code", shortCode))
	span.SetStatus(codes.Ok, "Short link created successfully")
}

// RedirectShortLink 处理短链接重定向的请求 (GET /{shortCode})
func (h *LinkHandler) RedirectShortLink(w http.ResponseWriter, r *http.Request, shortCode string) {
	ctx := r.Context()
	span := trace.SpanFromContext(ctx) // 获取由otelhttp中间件创建的span

	requestLogger := h.logger.With(
		slog.String("http_method", r.Method),
		slog.String("http_path", r.URL.Path),
		slog.String("short_code_param", shortCode),
	)
	requestLogger.InfoContext(ctx, "Handler: Received request to redirect short link.")
	span.AddEvent("Handler: Received redirect request", trace.WithAttributes(attribute.String("short_code", shortCode)))

	if strings.TrimSpace(shortCode) == "" {
		requestLogger.WarnContext(ctx, "Handler: Short code is empty in path.")
		span.SetStatus(codes.Error, "Short code is empty")
		http.NotFound(w, r) // 或者返回400 Bad Request
		return
	}

	// serviceCtx, serviceCancel := context.WithTimeout(ctx, 3*time.Second)
	// defer serviceCancel()
	serviceCtx := ctx

	linkEntry, err := h.svc.GetOriginalURL(serviceCtx, shortCode)
	if err != nil {
		span.RecordError(err)
		if errors.Is(err, store.ErrNotFound) {
			requestLogger.InfoContext(ctx, "Handler: Short code not found.", slog.String("short_code", shortCode))
			span.SetStatus(codes.Error, "Short code not found by service") // OTel规范中，NotFound通常也认为是Error状态
			http.NotFound(w, r)
		} else {
			requestLogger.ErrorContext(ctx, "Handler: Service failed to get original URL.", slog.Any("error", err), slog.String("short_code", shortCode))
			span.SetStatus(codes.Error, "Service error retrieving original URL")
			http.Error(w, "Error retrieving link.", http.StatusInternalServerError)
		}
		return
	}

	// 执行HTTP 302 Found重定向 (或 301 Moved Permanently 如果适用且希望客户端缓存)
	requestLogger.InfoContext(ctx, "Handler: Redirecting to original URL.",
		slog.String("short_code", shortCode),
		slog.String("original_url", linkEntry.OriginalURL),
	)
	span.SetAttributes(attribute.String("redirect_url", linkEntry.OriginalURL))
	span.SetStatus(codes.Ok, "Redirecting")
	http.Redirect(w, r, linkEntry.OriginalURL, http.StatusFound)
}
