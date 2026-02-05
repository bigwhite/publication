package main

import "fmt"

func IsValid(s string) bool {
	// 用切片模拟栈，存储 byte 或 rune
	stack := []byte{}
	// 建立右括号到左括号的映射，方便查找
	pairs := map[byte]byte{
		')': '(',
		']': '[',
		'}': '{',
	}

	for i := 0; i < len(s); i++ {
		char := s[i]

		// 如果是右括号
		if match, isRight := pairs[char]; isRight {
			// 1. 栈为空，说明右括号多了，非法
			if len(stack) == 0 {
				return false
			}
			// 2. 栈顶元素不匹配，非法
			top := stack[len(stack)-1]
			if top != match {
				return false
			}
			// 匹配成功，弹出栈顶
			stack = stack[:len(stack)-1]
		} else {
			// 如果是左括号，直接入栈
			stack = append(stack, char)
		}
	}

	// 最后栈必须为空，否则说明左括号多了
	return len(stack) == 0
}

func main() {
	cases := []string{"()", "()[]{}", "(]", "([)]", "{[]}"}
	for _, c := range cases {
		fmt.Printf("Input: %-8s Valid: %v\n", c, IsValid(c))
	}
}
