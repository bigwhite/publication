// go test -bench . -gcflags '-l'
package main

import "testing"

func BenchmarkCGO(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CallCFunc()
	}
}

func BenchmarkGo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CallGoFunc()
	}
}
