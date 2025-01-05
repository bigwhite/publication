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
	writer1 := os.Stdout
	writer2, err := os.Create("writer2.txt")
	if err != nil {
		fmt.Println("创建目标文件错误:", err)
		return
	}
	defer writer2.Close()
	multiWriter := io.MultiWriter(writer1, writer2)

	// multi reader -> multi writer
	n, err := io.Copy(multiWriter, multiReader)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("written %d bytes\n", n)
}
