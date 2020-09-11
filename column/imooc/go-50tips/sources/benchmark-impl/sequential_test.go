package bench

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	tls "github.com/huandu/go-tls"
)

var (
	m     map[int64]struct{} = make(map[int64]struct{}, 10)
	mu    sync.Mutex
	round int64 = 1
)

func BenchmarkSequential(b *testing.B) {
	fmt.Printf("\ngoroutine[%d] enter BenchmarkSequential: round[%d], b.N[%d]\n",
		tls.ID(), atomic.LoadInt64(&round), b.N)
	defer func() {
		atomic.AddInt64(&round, 1)
	}()

	for i := 0; i < b.N; i++ {
		mu.Lock()
		_, ok := m[round]
		if !ok {
			m[round] = struct{}{}
			fmt.Printf("goroutine[%d] enter loop in BenchmarkSequential: round[%d], b.N[%d]\n",
				tls.ID(), atomic.LoadInt64(&round), b.N)
		}
		mu.Unlock()
	}
	fmt.Printf("goroutine[%d] exit BenchmarkSequential: round[%d], b.N[%d]\n",
		tls.ID(), atomic.LoadInt64(&round), b.N)
}
