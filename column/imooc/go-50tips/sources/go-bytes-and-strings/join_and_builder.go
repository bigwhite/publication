package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	s := []string{"I", "love", "Go"}
	fmt.Println(strings.Join(s, " ")) // I love Go
	b := [][]byte{[]byte("I"), []byte("love"), []byte("Go")}
	fmt.Printf("%q\n", bytes.Join(b, []byte(" "))) // "I love Go"

	var builder strings.Builder
	for i, w := range s {
		builder.WriteString(w)
		if i != len(s)-1 {
			builder.WriteString(" ")
		}
	}
	fmt.Printf("%s\n", builder.String()) // I love Go

	var buf bytes.Buffer
	for i, w := range b {
		buf.Write(w)
		if i != len(b)-1 {
			buf.WriteString(" ")
		}
	}
	fmt.Printf("%s\n", buf.String()) // I love Go
}
