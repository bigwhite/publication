package main

import (
	"fmt"
)

// LogEntry 模拟一条日志结构
type LogEntry struct {
	Timestamp int64  // Unix 时间戳
	Content   string // 日志内容
}

// MergeLogStreams 模拟合并两个有序的日志流
// 这是一个经典的 "分离双指针" 模式应用
func MergeLogStreams(streamA, streamB []LogEntry) []LogEntry {
	// p1, p2 分别是 streamA 和 streamB 的读指针
	p1, p2 := 0, 0
	lenA, lenB := len(streamA), len(streamB)

	// 预分配结果切片，避免多次扩容带来的性能损耗
	// 这是一个 Go 工程实践的细节：Capacity Planning
	merged := make([]LogEntry, 0, lenA+lenB)

	// 循环条件：只要两个流中哪怕有一个还有数据，就要继续
	for p1 < lenA && p2 < lenB {
		// 比较两个指针当前指向的日志时间戳
		// 我们假设需要按时间升序排列（旧 -> 新）
		if streamA[p1].Timestamp < streamB[p2].Timestamp {
			merged = append(merged, streamA[p1])
			p1++ // 移动 A 的指针
		} else {
			merged = append(merged, streamB[p2])
			p2++ // 移动 B 的指针
		}
	}

	// 扫尾工作
	// 当一个流耗尽，另一个流可能还有剩余数据
	// 直接将剩余部分 append 进去即可，因为它们本身是有序的
	if p1 < lenA {
		merged = append(merged, streamA[p1:]...)
	}
	if p2 < lenB {
		merged = append(merged, streamB[p2:]...)
	}

	return merged
}

func main() {
	// 模拟数据：Stream A
	logsA := []LogEntry{
		{Timestamp: 100, Content: "[A] Server started"},
		{Timestamp: 103, Content: "[A] User login"},
		{Timestamp: 108, Content: "[A] Error occurred"},
	}

	// 模拟数据：Stream B
	logsB := []LogEntry{
		{Timestamp: 101, Content: "[B] Health check ok"},
		{Timestamp: 105, Content: "[B] Cache miss"},
		{Timestamp: 110, Content: "[B] Server shutdown"},
	}

	result := MergeLogStreams(logsA, logsB)

	fmt.Println("=== Merged Logs (Time Ordered) ===")
	for _, log := range result {
		fmt.Printf("[%d] %s\n", log.Timestamp, log.Content)
	}
}
