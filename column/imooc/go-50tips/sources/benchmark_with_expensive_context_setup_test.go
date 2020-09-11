package benchmark

import (
	"strings"
	"testing"
	"time"
)

var sl = []string{
	"Rob Pike ",
	"Robert Griesemer ",
	"Ken Thompson ",
}

func concatStringByJoin(sl []string) string {
	return strings.Join(sl, "")
}

func expensiveTestContextSetup() {
	time.Sleep(200 * time.Millisecond)
}

func BenchmarkStrcatWithTestContextSetup(b *testing.B) {
	expensiveTestContextSetup()
	for n := 0; n < b.N; n++ {
		concatStringByJoin(sl)
	}
}

func BenchmarkStrcatWithTestContextSetupAndResetTimer(b *testing.B) {
	expensiveTestContextSetup()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		concatStringByJoin(sl)
	}
}

func BenchmarkStrcatWithTestContextSetupAndRestartTimer(b *testing.B) {
	b.StopTimer()
	expensiveTestContextSetup()
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		concatStringByJoin(sl)
	}
}

func BenchmarkStrcat(b *testing.B) {
	for n := 0; n < b.N; n++ {
		concatStringByJoin(sl)
	}
}
