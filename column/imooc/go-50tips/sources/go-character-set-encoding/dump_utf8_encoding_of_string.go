package main

import "fmt"

func main() {
	var s = "中"
	fmt.Printf("Unicode字符：%s => 其UTF-8内存编码表示为: ", s)
	for _, v := range []byte(s) {
		fmt.Printf("0x%X ", v)
	}
	fmt.Printf("\n")
}
