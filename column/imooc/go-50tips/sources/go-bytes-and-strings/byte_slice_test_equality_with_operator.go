package main

import "fmt"

func main() {
	var a = []byte{'a', 'b', 'c'}
	var b = []byte{'a', 'b', 'd'}

	if a == b { // invalid operation: a == b
		fmt.Println("slice a is equal to slice b")
	} else {
		fmt.Println("slice a is not equal to slice b")
	}

	if a != nil { // valid operation
		fmt.Println("slice a is not nil")
	}
}
