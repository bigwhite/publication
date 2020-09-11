package main

import "fmt"

func Max(n int, m int) int {
	if n > m {
		return n
	} else {
		return m
	}
}

func main() {
	fmt.Printf("%d\n", Max(5, 6))
}
