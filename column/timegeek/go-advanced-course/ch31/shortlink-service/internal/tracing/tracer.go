// internal/tracing/tracer.go
package tracing

import (
	"fmt"
	"log/slog" // 使用slog记录初始化日志

	// "os" // 已经包含在slog或main中

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// InitTracerProvider initializes an OpenTelemetry TracerProvider with a stdout exporter.
// This is suitable for demos and local development to see traces in the console.
func InitTracerProvider(serviceName, serviceVersion string, enabled bool, sampleRatio float64) (*sdktrace.TracerProvider, error) {
	if !enabled {
		slog.Info("Distributed tracing is disabled by configuration.", slog.String("service_name", serviceName))
		// 如果禁用，返回nil，main函数中将不会设置全局TracerProvider，
		// otel.Tracer() 将返回一个NoOpTracer，不会产生实际的trace数据。
		return nil, nil
	}

	slog.Info("Initializing TracerProvider...",
		slog.String("service_name", serviceName),
		slog.String("exporter_type", "stdout"),
		slog.Float64("sample_ratio", sampleRatio),
	)

	// 1. 创建一个Exporter，这里使用标准输出 (stdouttrace)
	// WithPrettyPrint 使控制台输出的trace信息更易读。
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, fmt.Errorf("tracing: failed to create stdout trace exporter: %w", err)
	}
	slog.Debug("Stdout trace exporter initialized.")

	// 2. 定义资源 (Resource)，包含服务名、版本等通用属性
	// 这些属性会附加到所有由此Provider产生的Span上。
	res, err := resource.Merge(
		resource.Default(), // 包含默认属性如 telemetry.sdk.language, .name, .version
		resource.NewWithAttributes(
			semconv.SchemaURL, // OTel语义约定schema URL
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
			// 可以添加更多环境或部署相关的属性，例如：
			// attribute.String("deployment.environment", "development"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("tracing: failed to create OTel resource: %w", err)
	}
	slog.Debug("OTel resource defined.", slog.Any("resource_attributes", res.Attributes()))

	// 3. 配置采样器 (Sampler)
	var sampler sdktrace.Sampler
	if sampleRatio >= 1.0 {
		sampler = sdktrace.AlwaysSample() // 采样所有
	} else if sampleRatio <= 0.0 {
		sampler = sdktrace.NeverSample() // 不采样
	} else {
		// 根据给定的比例进行采样
		sampler = sdktrace.TraceIDRatioBased(sampleRatio)
	}
	// ParentBased确保如果上游服务（在分布式场景下）已经做出了采样决策，则遵循该决策；
	// 如果是根Trace（没有父Span），则使用我们上面配置的本地采样器(sampler)。
	finalSampler := sdktrace.ParentBased(sampler)
	slog.Debug("OTel sampler configured.", slog.Float64("effective_sample_ratio", sampleRatio))

	// 4. 创建TracerProvider，并配置SpanProcessor和Resource
	// NewBatchSpanProcessor将span批量异步导出，性能更好，是生产推荐（即使是对stdout exporter）。
	// NewSimpleSpanProcessor会同步导出每个span，仅用于非常简单的测试或调试。
	bsp := sdktrace.NewBatchSpanProcessor(exporter) // (可选) 配置批处理器参数，例如批处理超时、队列大小等
	// sdktrace.WithBatchTimeout(5*time.Second),
	// sdktrace.WithMaxQueueSize(2048),

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(finalSampler),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp), // 注册BatchSpanProcessor
	)
	slog.Debug("OTel TracerProvider created.")

	// 5. 设置为全局TracerProvider和全局TextMapPropagator
	// 这使得我们可以在应用的其他地方通过otel.Tracer("instrumentation-name")获取tracer实例，
	// 并通过otel.GetTextMapPropagator()进行上下文传播（主要用于分布式场景，但在单体内规范使用也有好处）。
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, // W3C Trace Context propagator (这是HTTP Headers的标准)
			propagation.Baggage{},      // W3C Baggage propagator (可选，用于传递业务数据)
		),
	)

	slog.Info("OpenTelemetry TracerProvider initialized and set globally with stdout exporter.", slog.String("service_name", serviceName))
	return tp, nil
}
