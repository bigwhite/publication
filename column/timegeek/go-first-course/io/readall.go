package main

import (
	"fmt"
	"io"
	"strings"
)

func main() {
	// 1. 正常读取
	reader := strings.NewReader("Hello, World!") // length = 13
	buf, err := io.ReadAll(reader)
	fmt.Printf("正常读取：Read %d bytes: %s, error=%v\n", len(buf), string(buf), err)

	// 2. 未读到任何数据
	reader1 := strings.NewReader("") // length = 0
	buf1, err := io.ReadAll(reader1)
	fmt.Printf("未读到任何数据：Read %d bytes: %s, error=%v\n", len(buf1), string(buf1), err)
}
