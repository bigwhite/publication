package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// ---- HTTP Server Metrics ----

// HTTPRequestsTotal 是一个CounterVec，用于记录HTTP请求的总数。
// 它按请求方法(method)、请求路径(path)和响应状态码(status_code)进行区分。
var HTTPRequestsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "shortlink", // 指标的命名空间，有助于组织和避免冲突
		Subsystem: "http_server",
		Name:      "requests_total", // 完整指标名将是 shortlink_http_server_requests_total
		Help:      "Total number of HTTP requests processed by the shortlink service.",
	},
	[]string{"method", "path", "status_code"}, // 标签名列表
)

// HTTPRequestDurationSeconds 是一个HistogramVec，用于观察HTTP请求延迟的分布情况。
var HTTPRequestDurationSeconds = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "shortlink",
		Subsystem: "http_server",
		Name:      "request_duration_seconds",
		Help:      "Histogram of HTTP request latencies for the shortlink service.",
		Buckets:   prometheus.DefBuckets, // prometheus.DefBuckets 是一组预定义的、通用的延迟桶
		// 或者，你可以根据你的服务特性自定义桶的边界，例如：
		// Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
	},
	[]string{"method", "path"}, // 按方法和路径区分
)

// ---- Application Specific Metrics (示例) ----

// ShortLinkCreationsTotal 记录成功创建的短链接总数
var ShortLinkCreationsTotal = promauto.NewCounter(
	prometheus.CounterOpts{
		Namespace: "shortlink",
		Subsystem: "service",
		Name:      "creations_total",
		Help:      "Total number of short links successfully created.",
	},
)

// ShortLinkRedirectsTotal 记录短链接重定向的总数，按状态（成功/未找到）区分
var ShortLinkRedirectsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "shortlink",
		Subsystem: "service",
		Name:      "redirects_total",
		Help:      "Total number of short link redirects, by status.",
	},
	[]string{"status"}, // "success", "not_found", "error"
)

// Init 初始化并注册所有必要的收集器。
// 这个函数应该在应用启动的早期被调用，例如在main.go中。
func Init() {
	// (可选) 注册构建信息指标 (go_build_info)
	prometheus.MustRegister(collectors.NewBuildInfoCollector())

	// 我们自定义的指标（如HTTPRequestsTotal, ShortLinkCreationsTotal等）
	// 因为使用了promauto包，它们在定义时已经自动注册到prometheus.DefaultRegisterer了，
	// 所以这里不需要再次显式调用 prometheus.MustRegister() 来注册它们。
	// 如果我们没有用promauto，而是用 prometheus.NewCounterVec() 等，则需要在这里注册。
	// 例如：
	// httpRequestsTotalPlain := prometheus.NewCounterVec(...)
	// prometheus.MustRegister(httpRequestsTotalPlain)
}
