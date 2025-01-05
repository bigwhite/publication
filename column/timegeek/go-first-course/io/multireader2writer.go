package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	reader1 := strings.NewReader("Hello, ")
	reader2 := strings.NewReader("World!\n")
	reader3, err := os.Open("reader3.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer reader3.Close()

	multiReader := io.MultiReader(reader1, reader2, reader3)

	// 将合并后的内容写到标准输出
	writer := os.Stdout

	_, err = io.Copy(writer, multiReader)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
