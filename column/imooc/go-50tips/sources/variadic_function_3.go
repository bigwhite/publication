package main

import "fmt"

func foo(b ...byte) {
	fmt.Println(string(b))
}

func main() {
	b := []byte{}
	b = append(b, "hello"...)
	fmt.Println(string(b))

	foo("hello"...)
}
