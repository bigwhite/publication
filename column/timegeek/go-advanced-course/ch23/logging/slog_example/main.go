package main

import (
	"log/slog"
	"os"
	"time"
)

type User struct {
	ID   string
	Name string
}

// UserLogValue implements slog.LogValuer to customize logging for User type
func (u User) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", u.ID),
		slog.String("name", u.Name),
	)
}

func main() {
	// --- TextHandler Example ---
	txtOpts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true, // 添加源码位置 (文件名:行号)
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// 将所有 Info 级别的 level 字段的键名改为 "severity"
			if a.Key == slog.LevelKey && a.Value.Any().(slog.Level) == slog.LevelInfo {
				a.Key = "severity"
			}
			return a
		},
	}
	txtHandler := slog.NewTextHandler(os.Stdout, txtOpts)
	txtLogger := slog.New(txtHandler)
	txtLogger.Info("TextHandler: Server started.", slog.String("port", ":8080"))

	// --- JSONHandler Example ---
	jsonOpts := &slog.HandlerOptions{
		Level: slog.LevelInfo, // JSON logger只输出Info及以上
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey { // 自定义时间戳字段名和格式
				a.Key = "event_time"
				a.Value = slog.StringValue(a.Value.Time().Format(time.RFC3339Nano))
			}
			return a
		},
	}
	jsonHandler := slog.NewJSONHandler(os.Stderr, jsonOpts)
	jsonLogger := slog.New(jsonHandler)
	jsonLogger.Error("JSONHandler: Payment failed.",
		slog.String("order_id", "ord-123"),
		slog.Any("error_details", map[string]string{"code": "P401", "message": "Insufficient funds"}),
	)

	// --- Using With for contextual logging ---
	requestID := "req-abc-789"
	userLogger := jsonLogger.With(
		slog.String("request_id", requestID),
		slog.Group("user_info",
			slog.String("id", "user-xyz"),
			slog.Bool("authenticated", true),
		),
	)
	userLogger.Info("User action performed.", slog.String("action", "view_profile"))
	userLogger.Debug("This debug from userLogger will not be printed by jsonLogger (LevelInfo).")

	// --- Implementing LogValuer for custom type logging ---
	currentUser := User{ID: "u-555", Name: "Alice Wonderland"}
	userLogger.Info("Processing user data.", slog.Any("user_object", currentUser))

	// --- Setting a default logger ---
	slog.SetDefault(txtLogger.WithGroup("global")) // All subsequent slog.Info etc. will use this
	slog.Info("This is a global log message via default logger (TextHandler).", slog.Int("count", 42))

	// Any package in your application can now simply call:
	// slog.Info("Message from another package")
	// without needing a logger instance passed around, if a default is set.
	// This is extremely convenient for widespread, simple logging needs.
	// However, for components requiring specific contexts or outputs, explicit logger injection is still preferred.
	doWorkWithGlobalLogger()
}

func doWorkWithGlobalLogger() {
	slog.Warn("Warning from a function using the global default logger.", slog.String("module", "worker"))
}
