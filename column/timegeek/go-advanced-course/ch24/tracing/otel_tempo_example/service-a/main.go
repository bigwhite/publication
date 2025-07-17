package main

import (
	"context"
	"demo/tracing" // 导入通用的tracing初始化包
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp" // HTTP client/server auto-instrumentation
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace" // 不直接用trace.Tracer，而是通过otel.Tracer获取
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serviceName := "service-a"
	serviceVersion := "1.0.0"
	otlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		otlpEndpoint = "localhost:14317" // OTel Collector服务地址
		log.Printf("[%s] OTEL_EXPORTER_OTLP_ENDPOINT not set, using default: %s\n", serviceName, otlpEndpoint)
	}

	// 初始化TracerProvider
	shutdownTracer, err := tracing.InitTracerProvider(ctx, serviceName, serviceVersion, otlpEndpoint)
	if err != nil {
		log.Fatalf("[%s] Failed to initialize TracerProvider: %v. Is OTel Collector running at %s?", serviceName, err, otlpEndpoint)
	}
	defer func() { // 确保在应用退出时关闭TracerProvider
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownTracer(shutdownCtx); err != nil {
			log.Printf("[%s] Error during TracerProvider shutdown: %v", serviceName, err)
		}
	}()

	// 获取一个Tracer实例
	tracer := otel.Tracer(serviceName + "-tracer") // Tracer命名

	// 创建一个带有OTel自动插桩的HTTP客户端
	// otelhttp.NewTransport 会自动为出站请求创建Span并注入Trace Context
	otelClient := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	// 定义HTTP Handler
	callBHandler := func(w http.ResponseWriter, r *http.Request) {
		// 从请求的context中启动一个新的Span，它会成为otelhttp.NewHandler创建的父Span的子Span
		// 或者如果这个handler是顶层入口，它会成为新的根Span（如果otelhttp.NewHandler没用）
		// 在本例中，我们将使用otelhttp.NewHandler包装整个Mux，所以这里tracer.Start会创建子Span
		requestCtx, parentSpan := tracer.Start(r.Context(), "service-a.handler.callServiceB")
		defer parentSpan.End()

		parentSpan.SetAttributes(attribute.String("http.target", r.URL.Path))
		log.Printf("[%s] Received request for %s\n", serviceName, r.URL.Path)

		// 获取service-b的URL (应来自配置或服务发现)
		serviceB_URL := os.Getenv("SERVICE_B_URL")
		if serviceB_URL == "" {
			serviceB_URL = "http://localhost:8081/data" // service-b服务地址
			log.Printf("[%s] SERVICE_B_URL not set, using default: %s\n", serviceName, serviceB_URL)
		}

		// 创建到service-b的请求，并使用带有当前Span的context
		// otelClient.Transport (otelhttp.NewTransport) 会自动从requestCtx中提取Trace Context并注入到出站请求头
		outboundReq, err := http.NewRequestWithContext(requestCtx, "GET", serviceB_URL, nil)
		if err != nil {
			parentSpan.RecordError(err)
			parentSpan.SetStatus(codes.Error, "failed to create request to service-b")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		log.Printf("[%s] Calling Service B at %s...\n", serviceName, serviceB_URL)
		resp, err := otelClient.Do(outboundReq) // 使用带OTel插桩的HTTP客户端发送请求
		if err != nil {
			parentSpan.RecordError(err)
			parentSpan.SetStatus(codes.Error, "failed to call service-b")
			http.Error(w, fmt.Sprintf("Failed to call service-b: %v", err), http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			parentSpan.RecordError(err)
			parentSpan.SetStatus(codes.Error, "failed to read response from service-b")
			http.Error(w, "Internal server error reading response", http.StatusInternalServerError)
			return
		}

		responseMessage := fmt.Sprintf("Service A got response from Service B: [%s]", string(bodyBytes))
		parentSpan.AddEvent("Received response from Service B", oteltrace.WithAttributes(attribute.Int("response.size", len(bodyBytes))))
		parentSpan.SetStatus(codes.Ok, "Successfully called service-b")

		w.WriteHeader(resp.StatusCode)
		fmt.Fprint(w, responseMessage)
		log.Printf("[%s] Successfully handled /call-b request.\n", serviceName)
	}

	// 使用otelhttp.NewHandler包装我们的业务handler，使其自动处理入站请求的Trace上下文和根Span创建
	// "service-a-http-server" 将作为这个HTTP服务器instrumentation的名称，影响根Span的命名
	tracedCallBHandler := otelhttp.NewHandler(http.HandlerFunc(callBHandler), "service-a.http.inbound")

	mux := http.NewServeMux()
	mux.Handle("/call-b", tracedCallBHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux, // 使用已包装的handler
	}

	go func() {
		log.Printf("[%s] HTTP server listening on :8080\n", serviceName)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[%s] Service A failed to start: %v", serviceName, err)
		}
	}()

	<-ctx.Done() // 等待退出信号
	log.Printf("[%s] Shutdown signal received, stopping server...\n", serviceName)
	shutdownServerCtx, cancelShutdownServer := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdownServer()
	if err := server.Shutdown(shutdownServerCtx); err != nil {
		log.Printf("[%s] Error during server shutdown: %v", serviceName, err)
	}
	log.Printf("[%s] Server stopped.\n", serviceName)
}
