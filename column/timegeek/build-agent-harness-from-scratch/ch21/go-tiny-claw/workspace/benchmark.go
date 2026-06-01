package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// 测试函数类型
type TestFunc func(int64, *sync.WaitGroup)

// 测试互斥锁方案
func testMutex(count int64, wg *sync.WaitGroup) {
	var mu sync.Mutex
	var c int64
	
	for i := 0; i < int(count); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			c++
			mu.Unlock()
		}()
	}
}

// 测试原子操作 AddInt64 方案
func testAtomicAdd(count int64, wg *sync.WaitGroup) {
	var c int64
	
	for i := 0; i < int(count); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&c, 1)
		}()
	}
}

// 测试 CAS 方案
func testCAS(count int64, wg *sync.WaitGroup) {
	var c int64
	
	for i := 0; i < int(count); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				current := atomic.LoadInt64(&c)
				if atomic.CompareAndSwapInt64(&c, current, current+1) {
					break
				}
			}
		}()
	}
}

// 运行性能测试
func runPerformanceTest(name string, testFunc TestFunc, count int64) {
	var wg sync.WaitGroup
	
	start := time.Now()
	testFunc(count, &wg)
	wg.Wait()
	duration := time.Since(start)
	
	fmt.Printf("%s - 耗时: %v\n", name, duration)
}

func main() {
	const count = int64(1000000) // 100万次操作
	
	fmt.Printf("开始性能测试，并发次数: %d\n", count)
	fmt.Println("=====================================")
	
	// 测试互斥锁
	runPerformanceTest("互斥锁方案", testMutex, count)
	
	// 测试原子操作 AddInt64
	runPerformanceTest("原子操作 AddInt64 方案", testAtomicAdd, count)
	
	// 测试 CAS 方案
	runPerformanceTest("CAS 方案", testCAS, count)
	
	fmt.Println("=====================================")
	fmt.Println("性能测试完成")
}