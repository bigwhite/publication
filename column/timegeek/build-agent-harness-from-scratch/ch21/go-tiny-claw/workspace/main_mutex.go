package main

import (
	"fmt"
	"sync"
)

func main() {
	// 全局计数器
	var count int
	var wg sync.WaitGroup
	var mu sync.Mutex // 互斥锁

	// 启动 1000 个 Goroutine 去并发累加
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()         // 加锁
			defer mu.Unlock() // 解锁
			count++
		}()
	}

	wg.Wait()
	fmt.Printf("最终的 Count 是: %d\n", count)
}