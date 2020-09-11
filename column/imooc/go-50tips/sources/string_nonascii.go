package main

import "fmt"

func main() {
	// 中文字符  Unicode CodePoint(码点)   UTF8编码
	//  中		U+4E2D			E4B8AD
	//  国		U+56FD			E59BBD
	//  欢		U+6B22			E6ACA2
	//  迎		U+8FCE			E8BF8E
	//  您		U+60A8			E682A8
	s := "中国欢迎您"
	rs := []rune(s)
	sl := []byte(s)
	for i, v := range rs {
		var utf8Bytes []byte
		for j := i * 3; j < (i+1)*3; j++ {
			utf8Bytes = append(utf8Bytes, sl[j])
		}
		fmt.Printf("%s => %X => %X\n", string(v), v, utf8Bytes)
	}
}
