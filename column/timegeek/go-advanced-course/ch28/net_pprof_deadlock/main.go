package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"net/http"
	_ "net/http/pprof" // 导入 pprof 包
)

func main() {
	runtime.SetMutexProfileFraction(1) // 开启 mutex分析
	runtime.SetBlockProfileRate(1)     // 开启 block 分析

	// 启动 pprof 监控
	go func() {
		log.Println("Starting pprof server at :6060")
		log.Println(http.ListenAndServe("localhost:6060", nil)) // pprof 端口
	}()

	fmt.Printf("My PID is: %d. Waiting for deadlock...\n", os.Getpid())
	var mu1, mu2 sync.Mutex

	var wg sync.WaitGroup
	wg.Add(2)

	go func() { // Goroutine 1
		defer wg.Done()
		mu1.Lock()
		fmt.Println("G1: mu1 locked")
		time.Sleep(100 * time.Millisecond) // Give G2 time to acquire mu2
		fmt.Println("G1: Attempting to lock mu2...")
		mu2.Lock() // Will block here waiting for G2
		fmt.Println("G1: mu2 locked (should not happen in deadlock)")
		mu2.Unlock()
		mu1.Unlock()
	}()

	go func() { // Goroutine 2
		defer wg.Done()
		mu2.Lock()
		fmt.Println("G2: mu2 locked")
		time.Sleep(100 * time.Millisecond) // Give G1 time to acquire mu1
		fmt.Println("G2: Attempting to lock mu1...")
		mu1.Lock() // Will block here waiting for G1
		fmt.Println("G2: mu1 locked (should not happen in deadlock)")
		mu1.Unlock()
		mu2.Unlock()
	}()

	var done = make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	for {
		select {
		case <-done:
			fmt.Println("Program close normally")
			return
		default:
			time.Sleep(5 * time.Second)
		}
	}
}
