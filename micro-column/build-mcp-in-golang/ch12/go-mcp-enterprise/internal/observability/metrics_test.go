package observability

import (
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestMetricsRecordAndHandler(t *testing.T) {
	metrics := NewMetrics()

	metrics.Record("get_order", 10*time.Millisecond, false)
	metrics.Record("get_order", 20*time.Millisecond, true)

	snapshot := metrics.Snapshot()
	got := snapshot["get_order"]

	if got.Calls != 2 {
		t.Fatalf("Calls = %d, want 2", got.Calls)
	}

	if got.Errors != 1 {
		t.Fatalf("Errors = %d, want 1", got.Errors)
	}

	if got.MaxDurationNS != uint64(20*time.Millisecond) {
		t.Fatalf(
			"MaxDurationNS = %d, want %d",
			got.MaxDurationNS,
			20*time.Millisecond,
		)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/metrics", nil)

	metrics.Handler().ServeHTTP(recorder, request)

	body := recorder.Body.String()

	if !strings.Contains(body, `mcp_tool_calls_total{tool="get_order"} 2`) {
		t.Fatalf("unexpected metrics body:\n%s", body)
	}

	if !strings.Contains(body, `mcp_tool_errors_total{tool="get_order"} 1`) {
		t.Fatalf("unexpected metrics body:\n%s", body)
	}
}

func TestMetricsConcurrentRecord(t *testing.T) {
	metrics := NewMetrics()

	const (
		goroutines = 20
		perWorker  = 100
	)

	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for j := 0; j < perWorker; j++ {
				metrics.Record("get_order", time.Millisecond, false)
			}
		}()
	}

	wg.Wait()

	got := metrics.Snapshot()["get_order"]
	want := uint64(goroutines * perWorker)

	if got.Calls != want {
		t.Fatalf("Calls = %d, want %d", got.Calls, want)
	}
}
