package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os" // For setting MutexProfileFraction etc.
	"time"

	"github.com/grafana/pyroscope-go"
)

// simulateSomeCPUWork performs some CPU-intensive work.
func simulateSomeCPUWork() {
	for i := 0; i < 100000000; i++ {
		_ = i * i
	}
}

// simulateMemoryAllocations allocates some memory.
func simulateMemoryAllocations() {
	for i := 0; i < 100; i++ {
		_ = make([]byte, 1024*1024) // Allocate 1MB
		time.Sleep(50 * time.Millisecond)
	}
}

func main() {
	// --- Pyroscope Configuration ---
	// 这些通常来自配置或环境变量
	pyroscopeServerAddress := os.Getenv("PYROSCOPE_SERVER_ADDRESS") // e.g., "http://pyroscope-server:4040"
	if pyroscopeServerAddress == "" {
		pyroscopeServerAddress = "http://localhost:4040" // Default for local demo
		log.Println("PYROSCOPE_SERVER_ADDRESS not set, using default:", pyroscopeServerAddress)
	}
	appName := "my-go-app.pyroscope-demo"

	// (可选) 开启Mutex和Block profiling, 对性能有一定影响, 按需开启
	// runtime.SetMutexProfileFraction(5) // Report 1 out of 5 mutex contention events
	// runtime.SetBlockProfileRate(5)     // Report 1 out of 5 block events (e.g. channel send/recv, select)

	// --- Start Pyroscope Profiler ---
	profiler, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: appName,
		ServerAddress:   pyroscopeServerAddress, // Pyroscope Server URL

		Logger: pyroscope.StandardLogger,
		// (可选) Tags to attach to all profiles
		Tags: map[string]string{"hostname": os.Getenv("HOSTNAME")},

		ProfileTypes: []pyroscope.ProfileType{
			// these profile types are enabled by default:
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,

			// these profile types are optional:
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
		// (可选) HTTP client for pyroscope agent
		// HTTPClient: &http.Client{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatalf("Failed to start Pyroscope profiler: %v. Is Pyroscope server running at %s?", err, pyroscopeServerAddress)
	}
	defer profiler.Stop()
	log.Printf("Pyroscope profiler started for app '%s', sending data to %s\n", appName, pyroscopeServerAddress)

	// --- Your Application Logic ---
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Tag a specific piece of work with pyroscope.TagWrapper
		// This will add "handler":"root" tag to profiles collected during this function's execution
		pyroscope.TagWrapper(r.Context(), pyroscope.Labels("handler", "root"), func(ctx context.Context) {
			fmt.Fprintf(w, "Hello from %s! Simulating work...\n", appName)
			simulateSomeCPUWork() // Simulate some CPU load
		})
	})

	go func() { // Simulate some background memory allocations
		log.Println("Starting background memory allocation simulation...")
		simulateMemoryAllocations()
		log.Println("Background memory allocation simulation finished.")
	}()

	port := "8080"
	log.Printf("HTTP server listening on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
