package main

import "fmt"

func main() {
	var b = []int{1, 2, 3, 4}
	fmt.Printf("slice b: len(%d), cap(%d), elements(%v)\n",
		len(b), cap(b), b)

	b1 := b[:2]
	fmt.Printf("slice b1: len(%d), cap(%d), elements(%v)\n",
		len(b1), cap(b1), b1)

	fmt.Println("\nappend 11 to b1:")
	b1 = append(b1, 11)
	fmt.Printf("slice b1: len(%d), cap(%d), elements(%v)\n",
		len(b1), cap(b1), b1)
	fmt.Printf("slice b: len(%d), cap(%d), elements(%v)\n",
		len(b), cap(b), b)
	fmt.Println("\nappend 22 to b1:")
	b1 = append(b1, 22)
	fmt.Printf("slice b1: len(%d), cap(%d), elements(%v)\n",
		len(b1), cap(b1), b1)
	fmt.Printf("slice b: len(%d), cap(%d), elements(%v)\n",
		len(b), cap(b), b)

	fmt.Println("\nappend 33 to b1:")
	b1 = append(b1, 33)
	fmt.Printf("slice b1: len(%d), cap(%d), elements(%v)\n",
		len(b1), cap(b1), b1)
	fmt.Printf("slice b: len(%d), cap(%d), elements(%v)\n",
		len(b), cap(b), b)

	b1[0] *= 100
	fmt.Println("\nb1[0] multiply 100:")
	fmt.Printf("slice b1: len(%d), cap(%d), elements(%v)\n",
		len(b1), cap(b1), b1)
	fmt.Printf("slice b: len(%d), cap(%d), elements(%v)\n",
		len(b), cap(b), b)
}
