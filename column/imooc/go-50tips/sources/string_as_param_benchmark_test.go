package main

import (
	"testing"
)

var s string = `Go, also known as Golang, is a statically typed, compiled programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson. Go is syntactically similar to C, but with memory safety, garbage collection, structural typing, and communicating sequential processes (CSP)-style concurrency.`

func handleString(s string) {
	_ = s + "hello"
}

func handlePtrToString(s *string) {
	_ = *s + "hello"

}

func BenchmarkHandleString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		handleString(s)
	}
}
func BenchmarkHandlePtrToString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		handlePtrToString(&s)
	}
}
