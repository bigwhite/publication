package main

import "fmt"

func doIteration(m map[int]int) {
	fmt.Printf("{ ")
	for k, v := range m {
		fmt.Printf("[%d, %d] ", k, v)
	}
	fmt.Printf("}\n")
}

func main() {
	m := map[int]int{
		1: 11,
		2: 12,
		3: 13,
	}

	for i := 0; i < 3; i++ {
		doIteration(m)
	}
}
