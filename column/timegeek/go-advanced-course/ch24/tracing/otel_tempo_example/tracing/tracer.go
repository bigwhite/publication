package tracing

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// InitTracerProvider initializes and registers an OTLP gRPC TracerProvider.
// It returns a shutdown function that should be called by the application on exit.
func InitTracerProvider(ctx context.Context, serviceName, serviceVersion, otlpEndpoint string) (func(context.Context) error, error) {
	log.Printf("Initializing TracerProvider for service '%s' (v%s), OTLP endpoint: '%s'\n", serviceName, serviceVersion, otlpEndpoint)

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
			// attribute.String("environment", "demo"), // 可选的其他全局属性
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTel resource: %w", err)
	}

	// 创建到OTLP Collector的gRPC连接
	// 在生产中，应该使用安全的凭证 (e.g., grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	// 并处理连接错误和重试
	connCtx, cancelConn := context.WithTimeout(ctx, 5*time.Second) // 连接超时
	defer cancelConn()
	conn, err := grpc.DialContext(connCtx, otlpEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // 仅用于演示
		grpc.WithBlock(), // 阻塞直到连接成功或超时
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to OTLP collector at '%s': %w", otlpEndpoint, err)
	}
	log.Printf("Successfully connected to OTLP collector at %s\n", otlpEndpoint)

	// 创建OTLP Trace Exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		// 尝试关闭连接，如果创建exporter失败
		if cerr := conn.Close(); cerr != nil {
			log.Printf("Warning: failed to close gRPC connection after exporter creation failed: %v", cerr)
		}
		return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}
	log.Println("OTLP trace exporter initialized.")

	// 创建BatchSpanProcessor，这是生产推荐的
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)

	// 创建TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // 为了演示，采样所有trace
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	// 设置为全局TracerProvider和Propagator
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, // W3C Trace Context (标准)
		propagation.Baggage{},      // W3C Baggage
	))

	log.Printf("Global TracerProvider and Propagator set for service '%s'.\n", serviceName)

	// 返回一个关闭函数，它会关闭TracerProvider和gRPC连接
	shutdownFunc := func(shutdownCtx context.Context) error {
		log.Printf("Attempting to shutdown TracerProvider for service '%s'...\n", serviceName)
		var errs []error
		if err := tp.Shutdown(shutdownCtx); err != nil {
			errs = append(errs, fmt.Errorf("TracerProvider shutdown error: %w", err))
			log.Printf("Error shutting down TracerProvider for %s: %v\n", serviceName, err)
		} else {
			log.Printf("TracerProvider for %s shut down successfully.\n", serviceName)
		}
		if err := conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("gRPC connection close error: %w", err))
			log.Printf("Error closing gRPC connection for %s: %v\n", serviceName, err)
		} else {
			log.Printf("gRPC connection for %s closed successfully.\n", serviceName)
		}
		if len(errs) > 0 {
			// 可以将多个错误合并返回，这里简单返回第一个
			return fmt.Errorf("shutdown for service %s encountered errors: %v", serviceName, errs)
		}
		return nil
	}

	return shutdownFunc, nil
}
