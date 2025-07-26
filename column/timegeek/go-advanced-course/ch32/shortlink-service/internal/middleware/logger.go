package middleware

import (
	"log/slog"
	"net/http"
	"time"
	// "go.opentelemetry.io/otel/trace" // 如果要从ctx获取traceID
)

// RequestLogger 是一个HTTP中间件，用于记录每个请求的访问日志。
// 它应该在Tracing中间件之后，以便能获取到TraceID。
func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			// 包装ResponseWriter以捕获状态码 (复用Metrics中间件中的或定义一个新的)
			// 为简单起见，我们假设Metrics中间件已经包装了，或者我们只记录请求开始和结束。
			// 如果需要状态码，需要像Metrics中间件那样包装。
			// 这里我们用一个简化的responseWriter，仅为演示。
			rw := newLoggingResponseWriter(w)

			// 从上下文中提取TraceID和SpanID (如果Tracing已集成)
			// span := trace.SpanFromContext(r.Context())
			// traceID := span.SpanContext().TraceID().String()
			// spanID := span.SpanContext().SpanID().String()
			// (为使本文件独立，暂时不引入otel/trace，假设这些会由slog的handler或上层添加)

			// 创建一个与请求绑定的logger实例
			// 在实际项目中，traceID等应从r.Context()中获取（由Tracing中间件注入）
			requestLogger := logger.With(
				slog.String("http_method", r.Method),
				slog.String("http_path", r.URL.Path),
				slog.String("http_proto", r.Proto),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				// slog.String("trace_id", traceID), // 如果获取了
			)

			requestLogger.Info("HTTP request received") // 请求开始

			// 执行链中的下一个处理器
			next.ServeHTTP(rw, r)

			duration := time.Since(startTime)
			statusCode := rw.statusCode // 从包装的writer获取状态码

			// 请求结束日志
			level := slog.LevelInfo
			if statusCode >= 500 {
				level = slog.LevelError
			} else if statusCode >= 400 {
				level = slog.LevelWarn
			}

			requestLogger.LogAttrs(r.Context(), level, "HTTP request completed",
				slog.Int("status_code", statusCode),
				slog.Duration("duration", duration),
				slog.Int64("response_size_bytes", rw.bytesWritten), // 假设responseWriter记录了大小
			)
		})
	}
}

// loggingResponseWriter 包装 http.ResponseWriter 以捕获状态码和响应大小
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int64
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(b)
	lrw.bytesWritten += int64(n)
	return n, err
}
