package basics

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	if Add(1, 2) != 3 {
		t.Error("1 + 2 should be 3") // 报告错误
	}
}

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ { // b.N 由测试框架动态调整
		Add(100, 200)
	}
}

func ExampleAdd() {
	sum := Add(5, 10)
	fmt.Println(sum)
	// Output: 15
}
