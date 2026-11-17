package observability

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Snapshot struct {
	Calls           uint64
	Errors          uint64
	TotalDurationNS uint64
	MaxDurationNS   uint64
}

type toolMetrics struct {
	calls           atomic.Uint64
	errors          atomic.Uint64
	totalDurationNS atomic.Uint64
	maxDurationNS   atomic.Uint64
}

// Metrics 收集 MCP Tool 的基础运行指标。
type Metrics struct {
	mu    sync.RWMutex
	tools map[string]*toolMetrics
}

func NewMetrics() *Metrics {
	return &Metrics{
		tools: make(map[string]*toolMetrics),
	}
}

func (m *Metrics) getToolMetrics(name string) *toolMetrics {
	m.mu.RLock()
	tm := m.tools[name]
	m.mu.RUnlock()

	if tm != nil {
		return tm
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if tm = m.tools[name]; tm == nil {
		tm = &toolMetrics{}
		m.tools[name] = tm
	}
	return tm
}

// Record 记录一次 Tool 调用。
func (m *Metrics) Record(name string, duration time.Duration, failed bool) {
	if m == nil {
		return
	}

	tm := m.getToolMetrics(name)

	tm.calls.Add(1)
	if failed {
		tm.errors.Add(1)
	}

	ns := uint64(duration.Nanoseconds())
	tm.totalDurationNS.Add(ns)

	for {
		old := tm.maxDurationNS.Load()
		if ns <= old || tm.maxDurationNS.CompareAndSwap(old, ns) {
			break
		}
	}
}

// Snapshot 返回当前指标快照。
func (m *Metrics) Snapshot() map[string]Snapshot {
	result := make(map[string]Snapshot)

	if m == nil {
		return result
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for name, tm := range m.tools {
		result[name] = Snapshot{
			Calls:           tm.calls.Load(),
			Errors:          tm.errors.Load(),
			TotalDurationNS: tm.totalDurationNS.Load(),
			MaxDurationNS:   tm.maxDurationNS.Load(),
		}
	}

	return result
}

// Handler 返回 Prometheus 文本格式的指标 HTTP Handler。
//
// 注意：本示例只提供 Counter 和累计耗时、最大耗时，
// 并未实现完整 Histogram 分桶。
func (m *Metrics) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set(
			"Content-Type",
			"text/plain; version=0.0.4; charset=utf-8",
		)

		snapshot := m.Snapshot()

		names := make([]string, 0, len(snapshot))
		for name := range snapshot {
			names = append(names, name)
		}
		sort.Strings(names)

		fmt.Fprintln(w, "# HELP mcp_tool_calls_total Total MCP tool calls.")
		fmt.Fprintln(w, "# TYPE mcp_tool_calls_total counter")

		for _, name := range names {
			fmt.Fprintf(
				w,
				"mcp_tool_calls_total{tool=%q} %d\n",
				name,
				snapshot[name].Calls,
			)
		}

		fmt.Fprintln(w, "# HELP mcp_tool_errors_total Total MCP tool errors.")
		fmt.Fprintln(w, "# TYPE mcp_tool_errors_total counter")

		for _, name := range names {
			fmt.Fprintf(
				w,
				"mcp_tool_errors_total{tool=%q} %d\n",
				name,
				snapshot[name].Errors,
			)
		}

		fmt.Fprintln(w, "# HELP mcp_tool_duration_seconds_sum Total tool duration in seconds.")
		fmt.Fprintln(w, "# TYPE mcp_tool_duration_seconds_sum counter")

		for _, name := range names {
			seconds := float64(snapshot[name].TotalDurationNS) / float64(time.Second)
			fmt.Fprintf(
				w,
				"mcp_tool_duration_seconds_sum{tool=%q} %f\n",
				name,
				seconds,
			)
		}

		fmt.Fprintln(w, "# HELP mcp_tool_duration_seconds_max Maximum tool duration in seconds.")
		fmt.Fprintln(w, "# TYPE mcp_tool_duration_seconds_max gauge")

		for _, name := range names {
			seconds := float64(snapshot[name].MaxDurationNS) / float64(time.Second)
			fmt.Fprintf(
				w,
				"mcp_tool_duration_seconds_max{tool=%q} %f\n",
				name,
				seconds,
			)
		}
	})
}

// escapeLabelValue 用于需要手工构造 Prometheus 标签时的转义。
// 本示例使用 fmt 的 %q 输出标签值，因此不直接调用此函数。
func escapeLabelValue(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, "\n", `\n`)
	value = strings.ReplaceAll(value, `"`, `\"`)
	return value
}

// RegisterTool 注册一个 Tool 的指标。
// 即使 Tool 尚未被调用，也会在 /metrics 中暴露零值。
func (m *Metrics) RegisterTool(name string) {
	if m == nil {
		return
	}
	m.getToolMetrics(name)
}
