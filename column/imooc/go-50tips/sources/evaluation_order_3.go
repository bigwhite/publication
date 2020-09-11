package main

import "fmt"

var (
	a    = c
	b, c = f()
	d    = 3
)

func f() (int, int) {
	d++
	return d, d + 1
}

func main() {
	fmt.Println(a, b, c, d)
}
