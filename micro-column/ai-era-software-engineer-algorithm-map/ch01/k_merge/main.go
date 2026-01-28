package main

import (
	"container/heap"
	"fmt"
)

// LogEntry 模拟一条日志结构
type LogEntry struct {
	Timestamp int64  // 排序键
	Content   string // 日志内容
	SourceID  int    // 标记这条日志来自哪个流（0, 1, 2...）
}

// LogHeap 是一个最小堆，存储 LogEntry
type LogHeap []LogEntry

func (h LogHeap) Len() int           { return len(h) }
func (h LogHeap) Less(i, j int) bool { return h[i].Timestamp < h[j].Timestamp } // 最小堆
func (h LogHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *LogHeap) Push(x interface{}) {
	*h = append(*h, x.(LogEntry))
}

func (h *LogHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// MergeKLogStreams 合并 K 个有序日志流
func MergeKLogStreams(streams [][]LogEntry) []LogEntry {
	k := len(streams)
	h := &LogHeap{}
	heap.Init(h)

	// pointers 数组记录每个流当前读取到的位置（索引）
	pointers := make([]int, k)

	// 1. 初始化：将每个流的第一个元素放入堆中
	for i := 0; i < k; i++ {
		if len(streams[i]) > 0 {
			// 注意：我们在日志对象中记录了它来自哪个流 (SourceID: i)
			// 这样当它被弹出时，我们知道该从哪个流补充下一个元素
			entry := streams[i][0]
			entry.SourceID = i
			heap.Push(h, entry)
			pointers[i]++ // 指针后移
		}
	}

	merged := make([]LogEntry, 0)

	// 2. 循环处理：弹出最小 -> 补充新元素
	for h.Len() > 0 {
		// 弹出堆顶（当前所有流中时间戳最小的日志）
		minLog := heap.Pop(h).(LogEntry)
		merged = append(merged, minLog)

		// 找到这条日志来自哪个流
		streamIdx := minLog.SourceID

		// 如果该流还有剩余元素，将下一个元素推入堆中
		if pointers[streamIdx] < len(streams[streamIdx]) {
			nextLog := streams[streamIdx][pointers[streamIdx]]
			nextLog.SourceID = streamIdx // 标记来源
			heap.Push(h, nextLog)
			pointers[streamIdx]++ // 指针后移
		}
	}

	return merged
}

func main() {
	// 模拟 3 个有序日志流
	stream1 := []LogEntry{{100, "S1_Log1", 0}, {105, "S1_Log2", 0}}
	stream2 := []LogEntry{{102, "S2_Log1", 0}, {108, "S2_Log2", 0}}
	stream3 := []LogEntry{{101, "S3_Log1", 0}, {103, "S3_Log2", 0}, {110, "S3_Log3", 0}}

	streams := [][]LogEntry{stream1, stream2, stream3}
	result := MergeKLogStreams(streams)

	fmt.Println("=== Merged 3 Streams ===")
	for _, log := range result {
		fmt.Printf("[%d] %s\n", log.Timestamp, log.Content)
	}
}
