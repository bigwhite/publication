package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	s := "大家好"
	fmt.Printf("字符串\"%s\"的长度为%d\n", s, len(s))
	fmt.Printf("字符串\"%s\"中的字符个数%d\n", s, utf8.RuneCountInString(s))

	s1 := "hello"
	fmt.Printf("s1[0] = %c\n", s1[0])
	fmt.Printf("s1[1] = %c\n", s1[1])
	// s[0] = 'j' // cannot assign to s[0]

	b := []byte(s1)
	b[0] = 'j'
	fmt.Printf("字符串s1 = %s\n", s1)
	fmt.Printf("切片b = %q\n", b)

	var s2 string
	fmt.Println(s2 == "")     // true
	fmt.Println(len(s2) == 0) // true
	//fmt.Println(s2 == nil)    // invalid operation: s2 == nil (mismatched types string and nil)

}
