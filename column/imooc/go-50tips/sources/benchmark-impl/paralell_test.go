package bench

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	tls "github.com/huandu/go-tls"
)

var (
	m     map[int64]int = make(map[int64]int, 20)
	mu    sync.Mutex
	round int64 = 1
)

func BenchmarkParalell(b *testing.B) {
	fmt.Printf("\ngoroutine[%d] enter BenchmarkParalell: round[%d], b.N[%d]\n",
		tls.ID(), atomic.LoadInt64(&round), b.N)
	defer func() {
		atomic.AddInt64(&round, 1)
	}()

	b.RunParallel(func(pb *testing.PB) {
		id := tls.ID()
		fmt.Printf("goroutine[%d] enter loop func in BenchmarkParalell: round[%d], b.N[%d]\n", tls.ID(), atomic.LoadInt64(&round), b.N)
		for pb.Next() {
			mu.Lock()
			_, ok := m[id]
			if !ok {
				m[id] = 1
			} else {
				m[id] = m[id] + 1
			}
			mu.Unlock()
		}

		mu.Lock()
		count := m[id]
		mu.Unlock()

		fmt.Printf("goroutine[%d] exit loop func in BenchmarkParalell: round[%d], loop[%d]\n", tls.ID(), atomic.LoadInt64(&round), count)
	})

	fmt.Printf("goroutine[%d] exit BenchmarkParalell: round[%d], b.N[%d]\n",
		tls.ID(), atomic.LoadInt64(&round), b.N)
}
