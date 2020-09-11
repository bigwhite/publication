package main

import "fmt"

func foo(a, b int) (x, y int) {
	defer func() {
		x = x * 5
		y = y * 10
	}()

	x = a + 5
	y = b + 6
	return
}

func main() {
	x, y := foo(1, 2)
	fmt.Println("x=", x, "y=", y)
}
