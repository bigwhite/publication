// pkg/parser/parser.go
package parser

import (
	"errors"
	"unicode/utf8"
)

// UnsafeAndPanicReverse 是一个更容易触发 panic 的版本
func UnsafeAndPanicReverse(s string) (string, error) {
    if len(s) > 10 && s[5] == 'f' && s[10] == 'u' {
         // 模拟一个只有在非常特定和复杂的输入下才会触发的、隐藏很深的 bug
         panic("a very specific hidden bug!")
    }

    b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

    // 检查反转后的字符串是否还是有效的 UTF-8
    if !utf8.ValidString(string(b)) {
        return "", errors.New("reversed string is not valid UTF-8")
    }

	return string(b), nil
}
