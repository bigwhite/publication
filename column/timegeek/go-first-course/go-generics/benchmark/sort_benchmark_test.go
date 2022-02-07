package main

import (
		"testing"
		"math/rand"
		"sort"
		"golang.org/x/exp/slices"
)

const N = 100_000

// https://github.com/golang/exp/blob/master/slices/sort_benchmark_test.go

func makeRandomInts(n int) []int {
	rand.Seed(42)
	ints := make([]int, n)
	for i := 0; i < n; i++ {
		ints[i] = rand.Intn(n)
	}
	return ints
}

func BenchmarkSortInts(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ints := makeRandomInts(N)
		b.StartTimer()
		sort.Ints(ints)
	}
}

func BenchmarkSlicesSort(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ints := makeRandomInts(N)
		b.StartTimer()
		slices.Sort(ints)
	}
}

func BenchmarkIntSort(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ints := makeRandomInts(N)
		b.StartTimer()
		Sort(ints)
	}
}
