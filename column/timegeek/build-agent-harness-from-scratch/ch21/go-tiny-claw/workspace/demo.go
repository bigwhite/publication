package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// 测试函数类型
type TestFunc func(int64, *sync.WaitGroup) int64

// 原始有问题的代码（用于对比）
func originalRace(count int64, wg *sync.WaitGroup) int64 {
	var c int64
	for i := 0; i < int(count); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c++ // 竞态条件
		}()
	}
	wg.Wait()
	return c
}

// 修复方案一：互斥锁
func mutexSolution(count int64, wg *sync.WaitGroup) int64 {
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
	wg.Wait()
	return c
}

// 修复方案二：原子操作
func atomicSolution(count int64, wg *sync.WaitGroup) int64 {
	var c int64
	for i := 0; i < int(count); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&c, 1)
		}()
	}
	wg.Wait()
	return c
}

// 修复方案三：原子操作 AddInt64
func atomicAddSolution(count int64, wg *sync.WaitGroup) int64 {
	var c int64
	for i := 0; i < int(count); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&c, 1)
		}()
	}
	wg.Wait()
	return c
}

// 修复方案四：CAS 操作
func cassolution(count int64, wg *sync.WaitGroup) int64 {
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
	wg.Wait()
	return c
}

// 运行并发测试
func runConcurrencyTest(name string, testFunc TestFunc, count int64) int64 {
	var wg sync.WaitGroup
	
	start := time.Now()
	result := testFunc(count, &wg)
	duration := time.Since(start)
	
	fmt.Printf("%s - 结果: %d, 期望: %d, 耗时: %v\n", 
		name, result, count, duration)
	
	return result
}

func main() {
	const count = int64(1000) // 1000次操作，便于观察结果
	
	fmt.Printf("并发安全测试，并发次数: %d\n", count)
	fmt.Println("=====================================")
	
	// 测试原始竞态条件（应该不准确）
	fmt.Println("原始竞态条件代码（应该不准确）:")
	runConcurrencyTest("原始竞态条件", originalRace, count)
	
	fmt.Println("\n修复方案测试:")
	
	// 测试互斥锁方案
	runConcurrencyTest("互斥锁方案", mutexSolution, count)
	
	// 测试原子操作方案
	runConcurrencyTest("原子操作方案", atomicSolution, count)
	
	// 测试原子操作 AddInt64 方案
	runConcurrencyTest("原子操作 AddInt64 方案", atomicAddSolution, count)
	
	// 测试 CAS 方案
	runConcurrencyTest("CAS 方案", cassolution, count)
	
	fmt.Println("=====================================")
	fmt.Println("测试完成")
	
	// 使用 race detector 测试
	fmt.Println("\n使用 race detector 测试:")
	fmt.Println("go run -race concurrency_test.go")
}