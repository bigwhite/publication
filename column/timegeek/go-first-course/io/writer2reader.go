package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type CapWriter struct {
	writer io.Writer
}

func NewCapWriter(w io.Writer) *CapWriter {
	return &CapWriter{writer: w}
}

func (cw *CapWriter) Write(p []byte) (n int, err error) {
	// 将写入的数据转换为大写
	upper := strings.ToUpper(string(p))

	// 将转换后的数据写入底层的 Writer
	return cw.writer.Write([]byte(upper))
}

type ReverseWriter struct {
	writer io.Writer
}

func NewReverseWriter(w io.Writer) *ReverseWriter {
	return &ReverseWriter{writer: w}
}

func (rw *ReverseWriter) Write(p []byte) (n int, err error) {
	// 将写入的数据进行翻转
	reversed := reverseString(string(p))

	// 将翻转后的数据写入底层的 Writer
	return rw.writer.Write([]byte(reversed))
}

func reverseString(s string) string {
	runes := []rune(s)
	n := len(runes)
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-1-i] = runes[n-1-i], runes[i]
	}
	return string(runes)
}

func main() {
	// 创建第一个 Pipe，获取对应的 Reader 和 Writer
	reader1, writer1 := io.Pipe()

	// 创建 CapWriter，将其与第一个 Pipe 返回的 Writer 整合
	capWriter := NewCapWriter(writer1)

	// 创建第二个 Pipe，获取对应的 Reader 和 Writer
	reader2, writer2 := io.Pipe()

	// 创建 ReverseWriter，将其与第二个 Pipe 返回的 Writer 整合
	reverseWriter := NewReverseWriter(writer2)

	// 启动一个 goroutine，向 CapWriter 写入数据
	go func() {
		defer writer1.Close()
		for i := 0; i < 5; i++ {
			data := []byte("hello, world!\n")
			_, err := capWriter.Write(data)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			time.Sleep(time.Second)
		}
	}()

	// 启动另一个 goroutine，从第二个 Pipe 返回的 Reader 中读取数据并处理
	go func() {
		defer writer2.Close()
		io.Copy(reverseWriter, reader1)
	}()

	// 从第二个 Pipe 返回的 Reader 中读取翻转后的数据并处理
	n, err := io.Copy(os.Stdout, reader2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("\nRead %d bytes\n", n)

	// 关闭 Reader
	err = reader2.Close()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
