package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug" // 用于在panic时打印堆栈并保持进程
	"time"
)

var shouldPanicImmediately = false // 控制是否在panic后立即退出，还是等待调试

func criticalOperation(input string) {
	if input == "trigger_panic" {
		log.Println("CRITICAL: About to perform an operation that will panic!")
		var data []int
		// 故意制造一个越界panic
		fmt.Println("Accessing out of bounds:", data[5]) // PANIC!
	}
	log.Printf("CRITICAL: Operation with input '%s' completed successfully (simulated).\n", input)
}

func oopsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("OOPS_HANDLER: Received request for /oops, preparing to panic...")

	// 在实际生产panic场景，程序会直接退出或被重启策略拉起。
	// 为了演示调试，我们可以在这里加入一个延迟或特定条件，
	// 使得在panic发生后，进程不会立即消失，给我们附加调试器的时间。
	// 或者，如果 'shouldPanicImmediately' 为 false，我们捕获panic，打印堆栈，然后死循环等待调试。
	if !shouldPanicImmediately {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("PANIC RECOVERED (for debugging): %v\n", r)
				log.Println("Stacktrace from recover:")
				debug.PrintStack() // 打印当前goroutine的堆栈

				// 为了让Delve有机会附加，这里让goroutine进入一个可控的等待状态
				// 在真实的线上panic（未recover或recover后os.Exit），进程会终止。
				// 如果是K8s等环境，Pod会被重启。
				// Core Dump是分析这种已终止进程panic的常用手段。
				// 这里我们模拟的是一个“卡住”而非立即退出的panic场景。
				log.Println("Process will now hang, waiting for debugger to attach to its PID...")
				for { // 无限循环，让进程保持存活
					time.Sleep(1 * time.Minute)
				}
			}
		}()
	}

	// 触发panic
	criticalOperation("trigger_panic") // 这个调用会panic

	// 这行不会被执行
	fmt.Fprintln(w, "If you see this, something went wrong with the panic trigger!")
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("PING_HANDLER: Received request for /ping")
	fmt.Fprintln(w, "PONG!")
}

func main() {
	// 从环境变量读取是否立即panic退出的配置
	if os.Getenv("PANIC_IMMEDIATELY") == "true" {
		shouldPanicImmediately = true
	}

	pid := os.Getpid()
	log.Printf("Starting HTTP server on :8080 (PID: %d)...", pid)
	log.Printf("  Normal endpoint: http://localhost:8080/ping")
	log.Printf("  Panic endpoint:  http://localhost:8080/oops")
	if !shouldPanicImmediately {
		log.Println("  NOTE: On /oops panic, this demo server will recover, print stack, and hang for debugging.")
		log.Println("        In a real production scenario without such a recover, the process would terminate.")
	}

	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/oops", oopsHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
