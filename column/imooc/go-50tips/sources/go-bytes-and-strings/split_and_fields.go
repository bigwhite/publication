package main

import (
	"bytes"
	"fmt"
	"strings"
)

func main() {
	// 使用Fields相关函数分割字符串
	fmt.Printf("%q\n", strings.Fields("go java python"))                         // ["go" "java" "python"]
	fmt.Printf("%q\n", strings.Fields("\tgo  \f \u0085 \u00a0 java \n\rpython")) // ["go" "java" "python"]
	fmt.Printf("%q\n", strings.Fields(" \t \n\r   "))                            // []

	splitFunc := func(r rune) bool {
		return r == rune('\n')
	}
	fmt.Printf("%q\n", strings.FieldsFunc("\tgo  \f \u0085 \u00a0 java \n\n\rpython", splitFunc)) // ["\tgo  \f \u0085 \u00a0 java " "\rpython"]

	// 使用Fields相关函数分割字节切片
	fmt.Printf("%q\n", bytes.Fields([]byte("go java python")))                                          // ["go" "java" "python"]
	fmt.Printf("%q\n", bytes.Fields([]byte("\tgo  \f \u0085 \u00a0 java \n\rpython")))                  // ["go" "java" "python"]
	fmt.Printf("%q\n", bytes.Fields([]byte(" \t \n\r   ")))                                             // []
	fmt.Printf("%q\n", bytes.FieldsFunc([]byte("\tgo  \f \u0085 \u00a0 java \n\n\rpython"), splitFunc)) // ["\tgo  \f \u0085 \u00a0 java " "\rpython"]

	// 使用Split相关函数分割字符串
	fmt.Printf("%q\n", strings.Split("a,b,c", ","))            // ["a" "b" "c"]
	fmt.Printf("%q\n", strings.Split("a,b,c", "b"))            // ["a," ",c"]
	fmt.Printf("%q\n", strings.Split("Go社区欢迎你", ""))           // ["G" "o" "社" "区" "欢" "迎" "你"]
	fmt.Printf("%q\n", strings.Split("abc", "de"))             // ["abc"]
	fmt.Printf("%q\n", strings.SplitN("a,b,c,d", ",", 2))      // ["a" "b,c,d"]
	fmt.Printf("%q\n", strings.SplitN("a,b,c,d", ",", 3))      // ["a" "b" "c,d"]
	fmt.Printf("%q\n", strings.SplitAfter("a,b,c,d", ","))     // ["a," "b," "c," "d"]
	fmt.Printf("%q\n", strings.SplitAfterN("a,b,c,d", ",", 2)) // ["a," "b,c,d"]

	// 使用Split相关函数分割字节切片
	fmt.Printf("%q\n", bytes.Split([]byte("a,b,c"), []byte(",")))            // ["a" "b" "c"]
	fmt.Printf("%q\n", bytes.Split([]byte("a,b,c"), []byte("b")))            // ["a," ",c"]
	fmt.Printf("%q\n", bytes.Split([]byte("Go社区欢迎你"), nil))                  // ["G" "o" "社" "区" "欢" "迎" "你"]
	fmt.Printf("%q\n", bytes.Split([]byte("abc"), []byte("de")))             // ["abc"]
	fmt.Printf("%q\n", bytes.SplitN([]byte("a,b,c,d"), []byte(","), 2))      // ["a" "b,c,d"]
	fmt.Printf("%q\n", bytes.SplitN([]byte("a,b,c,d"), []byte(","), 3))      // ["a" "b" "c,d"]
	fmt.Printf("%q\n", bytes.SplitAfter([]byte("a,b,c,d"), []byte(",")))     // ["a," "b," "c," "d"]
	fmt.Printf("%q\n", bytes.SplitAfterN([]byte("a,b,c,d"), []byte(","), 2)) // ["a," "b,c,d"]
}
