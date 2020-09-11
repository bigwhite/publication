package main

import (
	"bytes"
	"fmt"
)

func main() {
	var b bytes.Buffer
	b.Write([]byte("Effective Go"))
	fmt.Println(b.String())
}
