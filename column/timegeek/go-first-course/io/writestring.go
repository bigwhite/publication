package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	writer := os.Stdout
	data := "Hello, World!\n"
	n, err := io.WriteString(writer, data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Written %d bytes\n", n)
}
