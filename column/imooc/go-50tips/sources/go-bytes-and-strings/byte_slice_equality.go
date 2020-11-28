package main

import (
	"bytes"
	"fmt"
)

func main() {
	fmt.Println(bytes.Equal([]byte{'a', 'b', 'c'}, []byte{'a', 'b', 'd'})) // false
	fmt.Println(bytes.Equal([]byte{'a', 'b', 'c'}, []byte{'a', 'b', 'c'})) // true
	fmt.Println(bytes.Equal([]byte{'a', 'b', 'c'}, []byte{'b', 'a', 'c'})) // false
	fmt.Println(bytes.Equal([]byte{}, nil))                                // true
}
