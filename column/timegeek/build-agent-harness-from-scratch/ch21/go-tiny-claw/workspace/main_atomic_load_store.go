package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	// 全局计数器
	var count int64
	var wg sync.WaitGroup

	// 启动 1000 个 Goroutine 去并发累加
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// 使用原子操作读取和存储
			current := atomic.LoadInt64(&count)
			atomic.StoreInt64(&count, current+1)
		}()
	}

	wg.Wait()
	fmt.Printf("最终的 Count 是: %d\n", count)
}