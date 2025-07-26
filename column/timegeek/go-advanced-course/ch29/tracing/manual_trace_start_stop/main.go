package main

import (
	"fmt"
	"log"
	"os"
	"runtime/trace"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker %d: starting\n", id)
	time.Sleep(time.Duration(id*100) * time.Millisecond) // 模拟不同时长的任务
	fmt.Printf("Worker %d: finished\n", id)
}

func main() {
	// 1. 创建追踪输出文件
	traceFile := "manual_trace.out"
	f, err := os.Create(traceFile)
	if err != nil {
		log.Fatalf("Failed to create trace output file %s: %v", traceFile, err)
	}
	// 使用defer确保文件在main函数结束时关闭
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Failed to close trace file %s: %v", traceFile, err)
		}
	}()

	// 2. 启动追踪，将数据写入文件f
	if err := trace.Start(f); err != nil {
		log.Fatalf("Failed to start trace: %v", err)
	}
	// 3. 核心：使用defer确保trace.Stop()在main函数退出前被调用，
	//    这样所有缓冲的追踪数据才会被完整写入文件。
	defer trace.Stop()

	log.Println("Runtime tracing started. Executing some concurrent work...")

	var wg sync.WaitGroup
	numWorkers := 5
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}
	wg.Wait() // 等待所有worker完成

	log.Printf("All workers finished. Stopping trace. Trace data saved to %s\n", traceFile)
	fmt.Printf("\nTo analyze the trace, run:\ngo tool trace %s\n", traceFile)
}
