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
			// 使用 CompareAndSwap 循环 CAS 操作
			for {
				current := atomic.LoadInt64(&count)
				if atomic.CompareAndSwapInt64(&count, current, current+1) {
					break
				}
			}
		}()
	}

	wg.Wait()
	fmt.Printf("最终的 Count 是: %d\n", count)
}