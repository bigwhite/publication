package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	data := "Hello, World!\n"
	reader := strings.NewReader(data)

	writer1 := os.Stdout
	writer2, err := os.Create("writer2.txt")
	if err != nil {
		fmt.Println("创建目标文件错误:", err)
		return
	}
	defer writer2.Close()
	writer := io.MultiWriter(writer1, writer2)

	// reader -> multi writer
	n, err := io.Copy(writer, reader)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("written %d bytes\n", n)
}
