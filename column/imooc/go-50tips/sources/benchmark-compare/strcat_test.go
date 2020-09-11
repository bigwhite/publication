package main

import (
	"strings"
	"testing"
)

var sl = []string{
	"Rob Pike ",
	"Robert Griesemer ",
	"Ken Thompson ",
}

func Strcat(sl []string) string {
	return concatStringByJoin(sl)
	//return concatStringByOperator(sl)
}

func concatStringByOperator(sl []string) string {
	var s string
	for _, v := range sl {
		s += v
	}
	return s
}

func concatStringByJoin(sl []string) string {
	return strings.Join(sl, "")
}

func BenchmarkStrcat(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Strcat(sl)
	}
}
