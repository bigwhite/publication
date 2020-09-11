package main

import "fmt"

const (
	s = "string constant"
)

func main() {
	var s1 string = "string variable"

	fmt.Printf("%T\n", s)                          // string
	fmt.Printf("%T\n", s1)                         // string
	fmt.Printf("%T\n", "temporary string literal") // string
}
