package main

import (
	"strings"
	"testing"
)

const numStrings = 1000

func BenchmarkPlus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := ""
		for j := 0; j < numStrings; j++ {
			s += "hello"
		}
	}
}

func BenchmarkStringBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var sb strings.Builder
		for j := 0; j < numStrings; j++ {
			sb.WriteString("hello")
		}
		_ = sb.String() // 获取最终字符串，防止编译器优化
	}
}
