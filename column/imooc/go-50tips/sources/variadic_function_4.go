package main

import (
	"fmt"
	"strings"
)

func concat(a, b int) string {
	return fmt.Printf("%d %d", a, b)
}

func concat(x, y string) string {
	return x + " " + y
}

func concat(s []string) string {
	return strings.Join(s, " ")
}

func main() {
	println(concat(1, 2))
	println(concat("hello", "gopher"))
	println(concat([]string{"hello", "gopher", "!"}))
}
