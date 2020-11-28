package main

import (
	"fmt"
	"time"
)

func allocSlice(min, high int) []int {
	var b = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Printf("slice b: len(%d), cap(%d), elements(%v)\n",
		len(b), cap(b), b)

	go func() {
		time.Sleep(time.Second)
		fmt.Printf("slice b: len(%d), cap(%d), elements(%v)\n",
			len(b), cap(b), b)
	}()
	return b[min:high]
}

func main() {
	b1 := allocSlice(3, 7)
	fmt.Printf("slice b1: len(%d), cap(%d), elements(%v)\n",
		len(b1), cap(b1), b1)

	b2 := b1[:6]
	fmt.Printf("slice b2: len(%d), cap(%d), elements(%v)\n",
		len(b2), cap(b2), b2)

	b2[4] += 10
	b2[5] += 10

	time.Sleep(2 * time.Second)
}
