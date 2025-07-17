package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http" // For prometheus.NewGoCollector example
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "myapp_http_requests_total",
			Help: "Total number of HTTP requests processed by the application.",
		},
		[]string{"method", "path", "status_code"},
	)

	httpRequestDurationSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "myapp_http_request_duration_seconds",
			Help:    "Histogram of HTTP request latencies.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func handleHello(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	time.Sleep(time.Duration(rand.Intn(500)+50) * time.Millisecond)
	statusCode := http.StatusOK
	if rand.Intn(10) == 0 {
		statusCode = http.StatusInternalServerError
		w.WriteHeader(statusCode)
		fmt.Fprintf(w, "Oops! Something went wrong.")
	} else {
		w.WriteHeader(statusCode)
		fmt.Fprintf(w, "Hello from Go Metrics App!")
	}
	duration := time.Since(startTime).Seconds()

	httpRequestsTotal.With(prometheus.Labels{
		"method":      r.Method,
		"path":        r.URL.Path,
		"status_code": fmt.Sprintf("%d", statusCode),
	}).Inc()

	httpRequestDurationSeconds.With(prometheus.Labels{
		"method": r.Method,
		"path":   r.URL.Path,
	}).Observe(duration)

	log.Printf("%s %s - %d, duration: %.3fs", r.Method, r.URL.Path, statusCode, duration)
}

func main() {
	http.HandleFunc("/hello", handleHello)

	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())

	go func() {
		log.Println("Metrics server listening on :9091")
		if err := http.ListenAndServe(":9091", metricsMux); err != nil {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()

	log.Println("Application server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start application server: %v", err)
	}
}
