package main

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// 使用预设的 Production config 创建 logger (JSON, InfoLevel, Caller, Stacktrace on Error)
	// 也可以自定义Config:
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel), // 输出Debug及以上级别
		Development: false,
		Encoding:    "json", // json 或 console
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
			EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码器
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // Flushes any buffered log entries

	logger.Debug("This is a zap debug message.",
		zap.String("component", "auth"),
		zap.Int("user_id_count", 1005),
	)
	logger.Info("Zap logger initialized.",
		zap.String("url", "http://example.com"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
	logger.Warn("Potential issue detected.", zap.String("warning_code", "W001"))
	logger.Error("Operation failed.", zap.Error(fmt.Errorf("network connection refused")))
}
