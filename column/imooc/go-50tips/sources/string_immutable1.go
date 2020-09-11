package main

import (
	"fmt"
)

func main() {
	// original string
	var s string = "hello"
	fmt.Println("original string:", s)

	// reslice it and try to change the original string
	sl := []byte(s)
	sl[0] = 't'
	fmt.Println("slice:", string(sl))
	fmt.Println("after reslice, the original string is:", string(s))
}
