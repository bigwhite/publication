package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	writer := os.Stdout
	data := "Hello, World!\n"
	reader := strings.NewReader(data)

	// reader -> writer
	n, err := io.Copy(writer, reader)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("written %d bytes\n", n)
}
