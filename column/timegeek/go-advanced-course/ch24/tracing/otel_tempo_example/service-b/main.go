package main

import (
	"context"
	"demo/tracing"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// "math/rand" // For simulateWork if needed

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func simulateWork(ctx context.Context, duration time.Duration, operationName string) {
	// 获取当前context中的tracer，创建一个子span
	tracer := otel.Tracer("service-b-worker-tracer") // 可以用更具体的tracer name
	_, span := tracer.Start(ctx, operationName)
	defer span.End()

	span.SetAttributes(attribute.Int64("work.duration.ns", duration.Nanoseconds()))
	log.Printf("[Service B] Worker: Starting %s (will take %v)\n", operationName, duration)
	time.Sleep(duration)
	log.Printf("[Service B] Worker: Finished %s\n", operationName)
	span.AddEvent("Work simulation completed")
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	// otelhttp.NewHandler 已经为这个请求创建了一个服务器端Span，并将其放入r.Context()
	// 我们可以从r.Context()中获取当前的Span，或者直接用它来创建子Span
	ctx := r.Context()
	tracer := otel.Tracer("service-b-handler-tracer") // 获取tracer

	// 手动创建一个子span来表示这个handler内部的特定业务逻辑
	var handlerSpan oteltrace.Span // Using oteltrace alias from global import
	ctx, handlerSpan = tracer.Start(ctx, "service-b.handler.processData")
	defer handlerSpan.End()

	handlerSpan.SetAttributes(attribute.String("handler.message", "Service B processing /data request"))
	log.Printf("[Service B] Received request at /data. TraceID: %s\n", oteltrace.SpanFromContext(ctx).SpanContext().TraceID())

	// 模拟一些工作
	simulateWork(ctx, 50*time.Millisecond, "databaseQuery")
	simulateWork(ctx, 30*time.Millisecond, "externalAPICall")

	fmt.Fprintln(w, "Data from Service B (processed)")
	handlerSpan.AddEvent("Successfully returned data from Service B")
	handlerSpan.SetStatus(codes.Ok, "Data processed and returned")
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serviceName := "service-b"
	serviceVersion := "1.0.0"
	otlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		otlpEndpoint = "localhost:14317" // Otel Collector服务名
		log.Printf("[%s] OTEL_EXPORTER_OTLP_ENDPOINT not set, using default: %s\n", serviceName, otlpEndpoint)
	}

	shutdownTracer, err := tracing.InitTracerProvider(ctx, serviceName, serviceVersion, otlpEndpoint)
	if err != nil {
		log.Fatalf("[%s] Failed to initialize TracerProvider: %v", serviceName, err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownTracer(shutdownCtx); err != nil {
			log.Printf("[%s] Error during TracerProvider shutdown: %v", serviceName, err)
		}
	}()

	// 使用otelhttp.NewHandler来自动为HTTP请求创建span并处理上下文传播
	// "service-b.http.inbound" 将作为这个HTTP服务器instrumentation的名称
	handlerWithTracing := otelhttp.NewHandler(http.HandlerFunc(dataHandler), "service-b.http.inbound")

	mux := http.NewServeMux()
	mux.Handle("/data", handlerWithTracing) // 注册带追踪的handler

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	go func() {
		log.Printf("[%s] HTTP server listening on :8081\n", serviceName)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[%s] Service B failed to start: %v", serviceName, err)
		}
	}()

	<-ctx.Done()
	log.Printf("[%s] Shutdown signal received, stopping server...\n", serviceName)
	shutdownServerCtx, cancelShutdownServer := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdownServer()
	if err := server.Shutdown(shutdownServerCtx); err != nil {
		log.Printf("[%s] Error during server shutdown: %v", serviceName, err)
	}
	log.Printf("[%s] Server stopped.\n", serviceName)
}
