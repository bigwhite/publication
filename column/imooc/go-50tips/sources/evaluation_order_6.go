package main

import "fmt"

func example() {
	n0, n1 := 1, 2
	n0, n1 = n0+n1, n0
	fmt.Println(n0, n1)
}

func main() {
	example()
}
