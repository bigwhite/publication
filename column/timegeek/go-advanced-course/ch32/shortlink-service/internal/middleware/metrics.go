// internal/middleware/metrics.go
package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	appMetrics "github.com/your_org/shortlink/internal/metrics"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// 跳过诊断端点自身的指标记录
		if path == "/metrics" || strings.HasPrefix(path, "/debug/pprof") {
			next.ServeHTTP(w, r)
			return
		}

		startTime := time.Now()
		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r) // 执行链中的下一个handler
		duration := time.Since(startTime).Seconds()
		statusCodeStr := strconv.Itoa(rw.statusCode)

		// 路径规范化 (重要，但此处简化)
		// 真实项目中，应从路由匹配结果获取模板路径，例如 "/item/{id}"
		// if strings.HasPrefix(path, "/api/items/") { path = "/api/items/{id}" }

		appMetrics.HTTPRequestsTotal.WithLabelValues(r.Method, path, statusCodeStr).Inc()
		appMetrics.HTTPRequestDurationSeconds.WithLabelValues(r.Method, path).Observe(duration)
	})
}
