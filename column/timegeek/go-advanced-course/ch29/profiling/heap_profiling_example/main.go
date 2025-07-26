package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof" // 匿名导入以注册pprof HTTP handlers
	"runtime"          // 用于 GC 和 MemStats
	"strconv"
	"sync"
	"time"
)

var (
	// globalCache 模拟一个可能导致内存泄漏的缓存
	globalCache = make(map[string][]byte)
	cacheMutex  sync.Mutex                                        // 保护 globalCache 的并发访问
	randSrc     = rand.New(rand.NewSource(time.Now().UnixNano())) // 用于生成随机数据
)

// addToCache 持续向全局缓存中添加数据，模拟内存泄漏
func addToCache() {
	log.Println("Goroutine 'addToCache' started: continuously adding items to globalCache.")
	for i := 0; ; i++ {
		key := "cache_key_for_leak_simulation_" + strconv.Itoa(i)

		// 模拟不同大小的数据，平均约0.75KB (512 + 1024/2)
		dataSize := 512 + randSrc.Intn(512)
		data := make([]byte, dataSize)  // 分配[]byte
		for j := 0; j < dataSize; j++ { // 填充随机数据
			data[j] = byte(randSrc.Intn(256))
		}

		cacheMutex.Lock()
		globalCache[key] = data // 将数据存入全局map
		cacheMutex.Unlock()

		if i%5000 == 0 && i != 0 { // 每5000次打印一次日志，避免刷屏
			log.Printf("[addToCache] Added %d items to cache. Current cache size: %d items.\n", i+1, len(globalCache))
			// 主动触发一次GC，以便pprof heap能看到更真实的inuse数据（可选）
			runtime.GC()
		}
		time.Sleep(1 * time.Millisecond) // 控制添加速度，避免瞬间撑爆内存
	}
}

// frequentSmallAllocs 模拟高频小对象分配
func frequentSmallAllocs() {
	log.Println("Goroutine 'frequentSmallAllocs' started: frequently allocating small temporary objects.")
	for {
		// 模拟在请求处理或一些计算中创建临时对象
		for i := 0; i < 1000; i++ {
			// 分配一个小字符串（通常会在堆上，如果逃逸或足够大）
			_ = fmt.Sprintf("temp_string_data_%d_and_some_padding", i)
			// 分配一个小结构体 (如果它逃逸到堆)
			// type TempStruct struct { A int; B string }
			// _ = &TempStruct{A:i, B:"temp"}
		}
		time.Sleep(50 * time.Millisecond) // 每轮分配后稍作停顿
	}
}

// handleStats 提供一个简单的HTTP端点来查看当前缓存大小和内存统计
func handleStats(w http.ResponseWriter, r *http.Request) {
	cacheMutex.Lock()
	numItems := len(globalCache)
	cacheMutex.Unlock()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Fprintf(w, "--- Cache Stats ---\n")
	fmt.Fprintf(w, "Current cache items: %d\n\n", numItems)
	fmt.Fprintf(w, "--- Memory Stats (runtime.MemStats) ---\n")
	fmt.Fprintf(w, "Alloc (bytes allocated and not yet freed): %v MiB\n", m.Alloc/1024/1024)
	fmt.Fprintf(w, "TotalAlloc (bytes allocated since program start): %v MiB\n", m.TotalAlloc/1024/1024)
	fmt.Fprintf(w, "Sys (total bytes of memory obtained from OS): %v MiB\n", m.Sys/1024/1024)
	fmt.Fprintf(w, "NumGC: %v\n", m.NumGC)
	// ... 可以打印更多 MemStats 字段
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds) // 为日志添加微秒时间戳

	// 启动pprof HTTP服务器
	go func() {
		log.Println("Starting pprof HTTP server on localhost:6060")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("Pprof server failed: %v", err)
		}
	}()

	// 启动模拟内存泄漏的goroutine
	go addToCache()

	// 启动模拟高频小对象分配的goroutine
	go frequentSmallAllocs()

	// 启动业务HTTP服务器（用于查看stats和触发业务逻辑，如果未来添加）
	http.HandleFunc("/stats", handleStats)
	port := "8080"
	log.Printf("Business server listening on :%s. Access /stats to see cache size.\n", port)
	log.Println("Access http://localhost:6060/debug/pprof/heap to get heap profile.")
	log.Println("Let the application run for a while to observe memory changes.")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Business server failed: %v", err)
	}
}
