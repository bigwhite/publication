package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // 匿名导入pprof包
	"strconv"
	"sync"

	"time"
)

// calculateHashes: 一个CPU密集型任务，计算SHA256哈希多次
func calculateHashes(input string, iterations int) string {
	data := []byte(input)
	var hash [32]byte // sha256.Size
	for i := 0; i < iterations; i++ {
		hash = sha256.Sum256(data)
		// 为了让每次迭代的输入都不同，从而避免一些编译器优化或缓存效应，
		// 并且增加计算量，我们将上一次的哈希结果作为下一次的输入。
		// 注意：这只是为了模拟CPU消耗，实际意义不大。
		if i < iterations-1 { // 避免最后一次转换，因为我们只用hash
			data = hash[:]
		}
	}
	return fmt.Sprintf("%x", hash) // 返回最终哈希的十六进制字符串
}

// buildLongString: 另一个可能消耗CPU的函数，通过低效的字符串拼接
func buildLongString(count int) string {
	var s string // 使用+=进行字符串拼接，效率较低
	for i := 0; i < count; i++ {
		s += "Iteration " + strconv.Itoa(i) + " and some more text to make it longer. "
	}
	return s
}

// handleRequest: 模拟一个HTTP请求处理器，它会并发执行上述两个任务
func handleRequest(w http.ResponseWriter, r *http.Request) {
	iterations := 100000     // 哈希计算的迭代次数
	stringBuildCount := 2000 // 字符串拼接的迭代次数

	// 从查询参数中获取迭代次数，以便调整负载
	if queryIters := r.URL.Query().Get("iters"); queryIters != "" {
		if val, err := strconv.Atoi(queryIters); err == nil && val > 0 {
			iterations = val
		}
	}
	if queryStrCount := r.URL.Query().Get("strcount"); queryStrCount != "" {
		if val, err := strconv.Atoi(queryStrCount); err == nil && val > 0 {
			stringBuildCount = val
		}
	}

	log.Printf("Handling request: iterations=%d, stringBuildCount=%d\n", iterations, stringBuildCount)

	var wg sync.WaitGroup
	wg.Add(2) // 我们要等待两个goroutine完成

	var hashResult string
	var stringResultLength int // 只关心长度以避免打印过长字符串

	go func() { // Goroutine 1: 执行哈希计算
		defer wg.Done()
		startTime := time.Now()
		hashResult = calculateHashes("some_initial_seed_data_for_hashing", iterations)
		log.Printf("Hash calculation finished in %v. Result starts with: %s...\n",
			time.Since(startTime), hashResult[:min(10, len(hashResult))])
	}()

	go func() { // Goroutine 2: 执行字符串拼接
		defer wg.Done()
		startTime := time.Now()
		longStr := buildLongString(stringBuildCount)
		stringResultLength = len(longStr)
		log.Printf("String building finished in %v. Result length: %d\n",
			time.Since(startTime), stringResultLength)
	}()

	wg.Wait() // 等待两个goroutine都完成

	// 响应客户端
	fmt.Fprintf(w, "Work completed.\nHash result (prefix): %s...\nString result length: %d\n",
		hashResult[:min(10, len(hashResult))], stringResultLength)
}

// min是一个辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// 启动pprof HTTP服务器 (通常与业务服务器在同一端口，或一个单独的admin端口)
	go func() {
		log.Println("Starting pprof server on :6060")
		// http.ListenAndServe的第二个参数是handler，nil表示使用http.DefaultServeMux
		// _ "net/http/pprof" 会将其handlers注册到DefaultServeMux
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("Pprof server failed: %v", err)
		}
	}()

	// 启动业务HTTP服务器
	http.HandleFunc("/work", handleRequest) // 注册我们的业务处理器
	port := "8080"
	log.Printf("Business server listening on :%s. Access /work to generate load.\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil { // 使用DefaultServeMux
		log.Fatalf("Business server failed: %v", err)
	}
}
