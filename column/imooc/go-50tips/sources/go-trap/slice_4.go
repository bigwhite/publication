package main

import "fmt"

func allocSlice(min, high int) []int {
	var b = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Printf("slice b: len(%d), cap(%d), elements(%v)\n",
		len(b), cap(b), b)

	nb := make([]int, high-min, high-min)
	copy(nb, b[min:high])
	return nb
}

func main() {
	b1 := allocSlice(3, 7)
	fmt.Printf("slice b1: len(%d), cap(%d), elements(%v)\n",
		len(b1), cap(b1), b1)

	b2 := b1[:6]
	fmt.Printf("slice b2: len(%d), cap(%d), elements(%v)\n",
		len(b2), cap(b2), b2)
}
