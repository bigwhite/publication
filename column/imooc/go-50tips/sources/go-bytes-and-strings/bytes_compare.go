package main

import (
	"bytes"
	"fmt"
)

func main() {
	var a = []byte{'a', 'b', 'c'}
	var b = []byte{'a', 'b', 'd'}
	var c = []byte{} //empty slice
	var d []byte     //nil slice

	fmt.Println(bytes.Compare(a, b))     // -1 a < b
	fmt.Println(bytes.Compare(b, a))     // 1 b < a
	fmt.Println(bytes.Compare(c, d))     // 0
	fmt.Println(bytes.Compare(c, nil))   // 0
	fmt.Println(bytes.Compare(d, nil))   // 0
	fmt.Println(bytes.Compare(nil, nil)) // 0 nil == nil
}
