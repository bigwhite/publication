package main

import (
	"fmt"
	"sync"
	"time"
)

// Bucket 代表时间窗口中的一个小格子
type Bucket struct {
	count     int   // 该格子内的请求计数
	startTime int64 // 该格子记录的时间段起始点（毫秒时间戳）
}

// SlidingWindowLimiter 滑动窗口限流器
type SlidingWindowLimiter struct {
	mu          sync.Mutex
	windowSize  int64    // 整个窗口的大小（毫秒），例如 1000ms
	bucketTime  int64    // 单个格子代表的时间跨度（毫秒），例如 100ms
	buckets     []Bucket // 环形数组，存储每个格子的数据
	maxRequests int      // 限流阈值（在 windowSize 时间内允许的最大请求数）
}

// NewLimiter 创建一个限流器
// windowSize: 窗口总时间
// bucketCount: 将窗口切分成多少个格子（越多越平滑，但内存和计算开销越大）
// maxReq: 阈值
func NewLimiter(windowSize time.Duration, bucketCount int, maxReq int) *SlidingWindowLimiter {
	ms := windowSize.Milliseconds()
	return &SlidingWindowLimiter{
		windowSize:  ms,
		bucketTime:  ms / int64(bucketCount),
		buckets:     make([]Bucket, bucketCount),
		maxRequests: maxReq,
	}
}

// Allow 判断是否允许当前请求通过
func (l *SlidingWindowLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().UnixMilli()

	// 1. 定位：找到当前时间对应的格子下标（环形映射）
	// 例如：当前时间 12345ms，bucketTime 100ms，bucketCount 10
	// 12345 / 100 = 123 (总格子数)
	// 123 % 10 = 3 (当前落在数组下标 3 的位置)
	idx := int((now / l.bucketTime) % int64(len(l.buckets)))

	// 2. 对齐：计算当前格子的标准起始时间
	// 例如：12345ms 对应的格子起始时间是 12300ms
	currentBucketStart := now - (now % l.bucketTime)

	// 3. 懒惰重置（Lazy Reset）：
	// 检查该位置的格子是否是“上一轮”留下来的旧数据。
	// 如果格子的 startTime 不等于当前计算出的起始时间，说明这个格子里的数据
	// 是很久以前（至少一个窗口周期前）写入的，已经过期了，必须重置。
	if l.buckets[idx].startTime != currentBucketStart {
		l.buckets[idx] = Bucket{count: 0, startTime: currentBucketStart}
	}

	// 4. 统计：遍历所有格子，计算当前窗口内的请求总数
	totalCount := 0
	for _, b := range l.buckets {
		// 核心判断：只统计那些在 [now - windowSize, now] 时间范围内的格子
		// 如果一个格子的 startTime 太旧，说明它已经滑出窗口了，忽略它
		if now-b.startTime < l.windowSize {
			totalCount += b.count
		}
	}

	// 5. 决策：判断是否超限
	if totalCount >= l.maxRequests {
		fmt.Printf("Blocked! Total: %d >= Max: %d\n", totalCount, l.maxRequests)
		return false
	}

	// 6. 放行：当前格子计数+1
	l.buckets[idx].count++
	fmt.Printf("Allowed. Total: %d\n", totalCount+1)
	return true
}

func main() {
	// 初始化：1秒窗口，切分10个格子（每个100ms），限流阈值 5 QPS
	limiter := NewLimiter(1*time.Second, 10, 5)

	fmt.Println("--- Start Requests ---")

	// 模拟连续发送 10 次请求，每次间隔 100ms
	// 预期：前 5 次通过，第 6 次开始被限流
	for i := 0; i < 10; i++ {
		fmt.Printf("[%s] Req %d: ", time.Now().Format("15:04:05.000"), i+1)
		limiter.Allow()
		time.Sleep(100 * time.Millisecond)
	}

	// 休息 1 秒，等待窗口滑动过去，旧数据过期
	fmt.Println("\n--- Sleep 1 Second (Window Slides) ---")
	time.Sleep(1 * time.Second)

	// 再次发送请求，预期应该允许通过
	fmt.Printf("[%s] Req After Sleep: ", time.Now().Format("15:04:05.000"))
	limiter.Allow()
}
