package main

import (
	"fmt"
	"io" // Required for io.MultiWriter
	"log/slog"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
)

func main() {
	// 1. 配置 lumberjack.Logger 作为 io.Writer
	logFilePath := "./logs/myapp_rotated.log" // 日志文件路径
	os.MkdirAll("./logs", os.ModePerm)        // 确保logs目录存在

	logFileWriter := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    1,    // 每个日志文件的最大大小 (MB)
		MaxBackups: 3,    // 保留的旧日志文件的最大数量
		MaxAge:     7,    // 保留的旧日志文件的最大天数 (天)
		Compress:   true, // 是否压缩旧的日志文件 (使用gzip)
		LocalTime:  true, // 使用本地时间命名备份文件
	}

	// 2. 创建一个slog Handler
	// 使用 io.MultiWriter 可以同时将日志输出到文件和控制台（方便开发调试）
	handlerOptions := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true, // 添加源码位置
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(a.Value.Time().Format("2006-01-02 15:04:05.000"))
			}
			return a
		},
	}

	// 可以选择输出格式，例如JSON或Text
	// logger := slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, logFileWriter), handlerOptions))
	logger := slog.New(slog.NewTextHandler(io.MultiWriter(os.Stdout, logFileWriter), handlerOptions))

	slog.SetDefault(logger) // 可选：设置为默认logger

	slog.Info("Application starting with file rotation enabled via lumberjack.")
	slog.Debug("This is a debug message that should appear in both console and file.")

	// 模拟大量日志输出以触发轮转 (1MB大约需要较多条目)
	// 调整循环次数和日志内容大小以实际触发轮转
	// 一条典型slog日志（带时间、级别、源、消息、几个属性）可能在100-300字节左右
	// 1MB = 1024 * 1024 bytes. 大约需要 3500 - 10000条日志.
	for i := 0; i < 5000; i++ {
		slog.Info("This is a test log entry to demonstrate rotation.",
			slog.Int("entry_number", i),
			slog.String("data_payload", fmt.Sprintf("some_long_data_string_padding_%d_%s", i, time.Now().String())),
		)
		if i%1000 == 0 && i != 0 {
			slog.Warn("Milestone log entry reached.", slog.Int("progress_mark", i/1000))
			time.Sleep(100 * time.Millisecond) // 稍微暂停，让日志有机会写入
		}
	}

	slog.Info("Application finished logging. Check the './logs' directory for rotated files.")

	// lumberjack 在写入时自动处理轮转，通常不需要显式关闭或同步其自身。
	// 但如果你的 Handler 或 slog.Logger 有 Sync 方法 (slog.Logger 本身没有)，
	// 并且你希望在程序退出前确保所有缓冲都已写入，可以尝试调用。
	// 例如，如果你用的是 zap 的 slog Handler 桥接，zap.Logger 有 Sync()。
	// 对于标准 slog Handler，它们直接写入 io.Writer，由 Writer 负责缓冲。
	// logFileWriter.Close() 也不是必须的，除非你想显式关闭文件句柄。
}
