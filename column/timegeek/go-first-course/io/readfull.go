package main

import (
	"fmt"
	"io"
	"strings"
)

func main() {
	// 1. 正常读取
	reader := strings.NewReader("Hello, World!") // length = 13
	buf := make([]byte, 10)                      // 创建大小为10的缓冲区
	n, err := io.ReadFull(reader, buf)
	fmt.Printf("正常读取：Read %d bytes: %s, error=%v\n", n, string(buf), err)

	// 2. 恰好读完
	reader1 := strings.NewReader("Hello, World!") // length = 13
	buf1 := make([]byte, 13)
	n, err = io.ReadFull(reader1, buf1)
	fmt.Printf("恰好读完：Read %d bytes: %s, error=%v\n", n, string(buf1), err)

	// 3. 读取不足
	reader2 := strings.NewReader("Hello, World!") // length = 13
	buf2 := make([]byte, 15)                      // 创建大小为15的缓冲区
	n, err = io.ReadFull(reader2, buf2)
	fmt.Printf("读取不足：Read %d bytes: %s, error=%v\n", n, string(buf2), err)

	// 4. 未读到任何数据
	reader3 := strings.NewReader("") // length = 0
	buf3 := make([]byte, 15)         // 创建大小为15的缓冲区
	n, err = io.ReadFull(reader3, buf3)
	fmt.Printf("未读到任何数据：Read %d bytes: %s, error=%v\n", n, string(buf3), err)
}
