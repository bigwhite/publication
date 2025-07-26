// internal/app/app.go
package app

import (
	"context"
	"errors"
	"fmt"
	"io" // For store.Close
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	// 应用内部包
	"github.com/your_org/shortlink/internal/config"
	"github.com/your_org/shortlink/internal/handler"
	appMetrics "github.com/your_org/shortlink/internal/metrics"
	"github.com/your_org/shortlink/internal/middleware"
	"github.com/your_org/shortlink/internal/service"
	"github.com/your_org/shortlink/internal/store"        // Store接口
	"github.com/your_org/shortlink/internal/store/memory" // 内存存储实现
	appTracing "github.com/your_org/shortlink/internal/tracing"

	// 第三方库
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// App 封装了应用的所有依赖和启动关闭逻辑
type App struct {
	appName        string
	serviceVersion string
	logger         *slog.Logger
	cfg            *config.Config
	store          store.Store
	shortenerSvc   service.ShortenerService
	tracerProvider *sdktrace.TracerProvider
	httpServer     *http.Server
}

// New 创建一个新的App实例，完成所有依赖的初始化和注入
func New(cfg *config.Config, logger *slog.Logger, appName, version string) (*App, error) {
	app := &App{
		cfg:            cfg,
		logger:         logger.With("component", "appcore"), // App自身也可以有组件标识
		appName:        appName,
		serviceVersion: version,
	}

	app.logger.Info("Creating new application instance...",
		slog.String("appName", app.appName),
		slog.String("version", app.serviceVersion),
	)

	// --- [App.New - 初始化阶段 1] ---
	// 初始化可观测性组件 (Tracing & Metrics)
	if app.cfg.Tracing.Enabled {
		var errInitTracer error
		// 1a. 初始化 Tracing (OpenTelemetry)
		app.tracerProvider, errInitTracer = appTracing.InitTracerProvider(
			app.appName,
			app.serviceVersion,
			app.cfg.Tracing.Enabled,
			app.cfg.Tracing.SampleRatio,
		)
		if errInitTracer != nil {
			app.logger.Error("Failed to initialize TracerProvider", slog.Any("error", errInitTracer))
			return nil, fmt.Errorf("app.New: failed to init tracer provider: %w", errInitTracer)
		}
	}
	// 1b. 初始化 Metrics (Prometheus Go运行时等)
	appMetrics.Init()
	app.logger.Info("Prometheus Go runtime metrics collectors registered.")

	// --- [App.New - 初始化阶段 2] ---
	// 初始化核心业务依赖 (层级：Store -> Service -> Handler)
	app.logger.Debug("Initializing core dependencies...")
	// 2a. 初始化 Store 层
	switch strings.ToLower(app.cfg.Store.Type) {
	case "memory":
		app.store = memory.NewStore(app.logger.With("datastore", "memory"))
		app.logger.Info("Initialized in-memory store.")
	default:
		err := fmt.Errorf("unsupported store type from config: %s", app.cfg.Store.Type)
		app.logger.Error("Failed to initialize store", slog.Any("error", err))
		return nil, err
	}

	// 2b. 初始化 Service 层 (注入Store)
	app.shortenerSvc = service.NewShortenerService(app.store, app.logger.With("layer", "service"), nil)
	app.logger.Info("Shortener service initialized.")

	// 2c. 初始化 Handler 层 (注入Service)
	linkHdlr := handler.NewLinkHandler(app.shortenerSvc, app.logger.With("layer", "handler"))
	app.logger.Info("Link handler initialized.")

	// --- [App.New - 初始化阶段 3] ---
	// 创建HTTP Router并注册所有路由
	mux := http.NewServeMux()
	// 3a. 注册业务路由
	mux.HandleFunc("POST /api/links", linkHdlr.CreateShortLink)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/metrics" || strings.HasPrefix(path, "/debug/pprof") {
			return // 这些由诊断路由处理
		}
		if r.Method == http.MethodGet && len(path) > 1 && path[0] == '/' && !strings.Contains(path[1:], "/") {
			shortCode := path[1:]
			linkHdlr.RedirectShortLink(w, r, shortCode)
			return
		}
		http.NotFound(w, r)
	})
	// 3b. 注册诊断路由 (Metrics, pprof)
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/debug/pprof/", http.DefaultServeMux.ServeHTTP) // 假设pprof已注册到DefaultServeMux

	app.logger.Info("HTTP routes registered.")

	// --- [App.New - 初始化阶段 4] ---
	// 应用HTTP中间件 (顺序很重要)
	var finalHandler http.Handler = mux
	// 4a. 应用 Metrics 中间件
	finalHandler = middleware.Metrics(finalHandler)
	app.logger.Info("Applied HTTP Metrics middleware.")
	// 4b. 应用 Tracing 中间件
	if app.tracerProvider != nil {
		finalHandler = otelhttp.NewHandler(finalHandler, fmt.Sprintf("%s.http.server", app.appName),
			otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
		)
		app.logger.Info("Applied OpenTelemetry HTTP Tracing middleware.")
	}
	// 4c. 应用 Logging 中间件
	finalHandler = middleware.RequestLogger(app.logger)(finalHandler)
	app.logger.Info("Applied HTTP Request Logging middleware.")

	// --- [App.New - 初始化阶段 5] ---
	// 创建并配置最终的HTTP服务器
	app.httpServer = &http.Server{
		Addr:         ":" + app.cfg.Server.Port,
		Handler:      finalHandler,
		ReadTimeout:  app.cfg.Server.ReadTimeout,
		WriteTimeout: app.cfg.Server.WriteTimeout,
		IdleTimeout:  app.cfg.Server.IdleTimeout,
	}
	app.logger.Info("HTTP server and dependencies initialized successfully.", slog.String("listen_addr", app.httpServer.Addr))
	return app, nil
}

// Run 启动应用并阻塞，直到接收到退出信号并完成优雅关闭
func (a *App) Run() error {
	a.logger.Info("Starting application run cycle...",
		slog.String("appName", a.appName),
		slog.String("version", a.serviceVersion),
	)

	// --- [App.Run - 运行阶段 1] ---
	// 异步启动HTTP服务器
	errChan := make(chan error, 1)
	go func() {
		a.logger.Info("HTTP server starting to listen and serve...")
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("HTTP server ListenAndServe failed", slog.Any("error", err))
			errChan <- err
		}
		close(errChan)
	}()
	a.logger.Info("HTTP server startup process initiated.")

	// --- [App.Run - 运行阶段 2] ---
	// 实现优雅退出：阻塞等待OS信号或服务器启动错误
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quitChannel:
		a.logger.Warn("Received shutdown signal, initiating graceful shutdown...",
			slog.String("signal", sig.String()),
		)
	case err := <-errChan:
		if err != nil {
			a.logger.Error("HTTP server failed to start, initiating shutdown...", slog.Any("error", err))
		}
	}

	// --- [App.Run - 运行阶段 3] ---
	// 执行优雅关闭流程
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), a.cfg.Server.ShutdownTimeout)
	defer cancelShutdown()

	a.logger.Info("Attempting to gracefully shut down the HTTP server...",
		slog.Duration("shutdown_timeout", a.cfg.Server.ShutdownTimeout),
	)
	if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
		a.logger.Error("HTTP server graceful shutdown failed or timed out", slog.Any("error", err))
	} else {
		a.logger.Info("HTTP server stopped gracefully.")
	}

	// --- [App.Run - 运行阶段 4] ---
	// 清理其他应用资源
	a.logger.Info("Cleaning up other application resources...")
	if a.store != nil {
		if closer, ok := a.store.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				a.logger.Error("Error closing store", slog.Any("error", err))
			} else {
				a.logger.Info("Store closed successfully.")
			}
		}
	}
	if a.tracerProvider != nil {
		tpShutdownCtx, tpCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer tpCancel()
		a.logger.Info("Attempting to shut down TracerProvider...")
		if err := a.tracerProvider.Shutdown(tpShutdownCtx); err != nil {
			a.logger.Error("Error shutting down TracerProvider", slog.Any("error", err))
		} else {
			a.logger.Info("TracerProvider shut down successfully.")
		}
	}

	a.logger.Info("Application has shut down completely.")
	return nil // 优雅关闭完成，返回nil
}

// Logger 返回App内部的logger实例，方便main函数在Run()返回错误时使用
func (a *App) Logger() *slog.Logger {
	return a.logger
}
