package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	// 创建一个新的logrus logger实例
	log := logrus.New()

	// 设置输出到标准输出 (logrus默认输出到stderr)
	log.SetOutput(os.Stdout)

	// 设置日志级别
	log.SetLevel(logrus.InfoLevel)

	// 设置JSON格式化器
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})

	// 添加固定的字段到所有日志条目
	log = log.WithFields(logrus.Fields{
		"service": "my-logrus-app",
		"version": "1.0.2",
	}).Logger // For new logrus versions, WithFields returns an Entry, use .Logger to get a Logger

	log.Info("A group of walrus emerges from the ocean")
	log.Warn("The group is larger than expected")
	log.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
		"action": "swim",
	}).Error("Failed to count all walruses")

	// Debug log (不会输出，因为级别是Info)
	log.Debug("This is a debug message for walrus counting.")
}
