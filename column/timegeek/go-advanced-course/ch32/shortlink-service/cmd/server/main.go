// cmd/server/main.go
package main

import (
	"fmt"
	"log/slog"         // 用于在app.New失败时记录日志
	_ "net/http/pprof" // 匿名导入pprof包，自动注册handlers到http.DefaultServeMux
	"os"
	"strings" // 用于initSlogLogger
	"time"    // 用于initSlogLogger

	"github.com/your_org/shortlink/internal/app"    // 导入app核心包
	"github.com/your_org/shortlink/internal/config" // 导入配置包
)

const (
	// defaultAppName 和 defaultAppVersion 作为硬编码的默认值，
	// 它们可以在没有配置文件或配置项时的最终兜底。
	// 在实际项目中，版本号更推荐通过编译时注入（ldflags）来设置。
	defaultAppName    = "ShortLinkService"
	defaultAppVersion = "0.1.0"
)

// bootstrapLogger 是一个引导日志器，专门用于应用启动的极早期阶段。
// 它的主要目的是在`app.New()`函数执行过程中或失败时，能够以结构化的方式记录日志。
// 它直接从传入的（可能已加载的）配置中获取日志级别和格式，
// 如果配置未加载，则使用安全的默认值。
// 注意：这个logger是临时的，App实例创建成功后，会使用其内部更完善的、基于完整配置的logger。
func bootstrapLogger(cfg *config.Config) *slog.Logger {
	var level slog.Level
	logLevelStr := "info"  // 安全的默认日志级别
	logFormatStr := "text" // 默认使用text格式，便于在启动时直接在控制台阅读

	// 如果配置文件已成功加载，则使用其中的设置
	if cfg != nil {
		logLevelStr = cfg.Server.LogLevel
		logFormatStr = cfg.Server.LogFormat
	}

	switch strings.ToLower(logLevelStr) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		// 如果配置中的日志级别无效，也使用安全的默认值
		level = slog.LevelInfo
	}

	var handler slog.Handler
	// 创建一个Handler，配置其日志级别和时间戳格式。
	// 注意：AddSource设为false，因为引导日志通常不需要源码位置，且可以提升一点点性能。
	// 日志输出到标准错误流(os.Stderr)，这是记录启动过程错误的常见做法。
	handlerOpts := &slog.HandlerOptions{AddSource: false, Level: level, ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			a.Value = slog.StringValue(a.Value.Time().Format(time.RFC3339))
		}
		return a
	}}

	if strings.ToLower(logFormatStr) == "json" {
		handler = slog.NewJSONHandler(os.Stderr, handlerOpts)
	} else {
		handler = slog.NewTextHandler(os.Stderr, handlerOpts)
	}

	// 返回一个带有"bootstrap_phase"标识的logger，便于区分这是启动阶段的日志。
	return slog.New(handler).With(slog.String("bootstrap_phase", "main"))
}

// main 函数是整个应用程序的入口点。
// 它的职责被设计得非常“薄”，主要负责引导和协调应用的创建与运行。
func main() {
	// --- 步骤 0: (可选) 解析最顶层的命令行参数 ---
	// 例如，通过`flag`包解析 `-config` 参数来指定配置文件的路径。
	// configFile := flag.String("config", "./configs/config.yaml", "Path to config file")
	// flag.Parse()
	// 为了保持本示例的简洁性，我们暂时直接使用固定的路径和名称。

	// --- 步骤 1: 加载配置 ---
	// 这是应用启动的第一步关键操作。后续所有组件的初始化都依赖于这份配置。
	// LoadConfig内部会处理文件查找、读取、解析以及环境变量的覆盖。
	cfg, err := config.LoadConfig("./configs", "config", "yaml")
	if err != nil {
		// 如果连配置都加载失败，这通常是一个致命错误，无法继续。
		// 此时日志系统可能还未初始化，所以使用 `fmt.Fprintf` 直接向标准错误输出。
		fmt.Fprintf(os.Stderr, "FATAL: Error loading initial application configuration: %v\n", err)
		os.Exit(1) // 以非零状态码退出，表示启动失败。
	}

	// --- 步骤 2: 初始化引导日志器 ---
	// 基于刚刚加载的配置（或其默认值），创建一个临时的引导日志器。
	// 这个logger主要用于记录`app.New()`的执行过程，以及在`app.New()`失败时能输出结构化的错误信息。
	mainLogger := bootstrapLogger(&cfg)
	// 我们可以在这里选择是否立即将此logger设置为全局默认实例(`slog.SetDefault`)。
	// 通常，更推荐的做法是在app.New()内部，当所有配置都最终确定后，
	// 创建并设置最终的、功能完备的全局默认logger。

	mainLogger.Info("Configuration loaded, attempting to create application instance...")

	// --- 步骤 3: 创建App实例 ---
	// 这是整个初始化过程的核心。我们将配置和引导日志器注入到`app.New()`函数中。
	// `app.New()`内部会完成所有组件（如Store, Service, Handler, HTTP Server等）的实例化和依赖注入。
	appNameFromConfig := cfg.AppName
	if appNameFromConfig == "" { // 如果配置中未指定appName，则使用硬编码的默认值
		appNameFromConfig = defaultAppName
	}

	application, err := app.New(&cfg, mainLogger, appNameFromConfig, defaultAppVersion)
	if err != nil {
		// 如果`app.New()`返回错误，说明应用的核心组件未能成功创建。
		// 我们使用刚刚创建的引导日志器来记录这个致命错误，然后退出。
		mainLogger.Error("FATAL: Failed to create application instance", slog.Any("error", err))
		os.Exit(1)
	}
	mainLogger.Info("Application instance created successfully.")

	// --- 步骤 4: 运行App ---
	// `application.Run()`是一个阻塞调用，它封装了应用的整个运行时生命周期，
	// 包括启动所有服务（如HTTP服务器）、监听操作系统信号、以及执行优雅关闭。
	if err := application.Run(); err != nil {
		// 如果`Run()`方法返回错误，表明应用在运行或关闭过程中遇到了未能处理的严重问题。
		// 此时，应用内部的、功能完备的日志系统应该已经可用，我们可以通过它来记录最终的致命错误。
		// (假设App实例有一个Logger()方法可以暴露其内部的logger)
		application.Logger().Error("FATAL: Application run failed and terminated", slog.Any("error", err))
		os.Exit(1)
	}

	// 如果`Run()`正常返回nil，说明应用已成功完成优雅关闭流程。
	// 相关的成功日志应在app.Run()内部或组件的Stop方法中打印。
	// main函数在此处正常结束，隐式返回os.Exit(0)。
	mainLogger.Info("Application main function exiting cleanly.")
}
