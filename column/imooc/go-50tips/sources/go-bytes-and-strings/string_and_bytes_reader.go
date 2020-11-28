package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func main() {
	var buf bytes.Buffer
	var s = "I love Go!!"

	_, err := io.Copy(&buf, strings.NewReader(s))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%q\n", buf.String()) // "I love Go!!"

	buf.Reset()
	var b = []byte("I love Go!!")
	_, err = io.Copy(&buf, bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%q\n", buf.String()) // "I love Go!!"
}
