package main

import "fmt"

func ValidateJSONStructure(jsonStr string) bool {
	stack := []byte{}
	inString := false  // 是否在字符串内部
	isEscaped := false // 上一个字符是否是转义符 \

	for i := 0; i < len(jsonStr); i++ {
		char := jsonStr[i]

		// 1. 处理字符串模式
		if inString {
			if isEscaped {
				// 如果上一个是转义符，当前字符无论是什么，都只是字符串的一部分
				// 重置转义状态
				isEscaped = false
			} else if char == '\\' {
				// 遇到转义符，标记一下，下一个字符将被“保护”
				isEscaped = true
			} else if char == '"' {
				// 遇到未转义的引号，说明字符串结束
				inString = false
			}
			// 在字符串里，忽略所有其他字符（包括 { } [ ]）
			continue
		}

		// 2. 处理普通模式
		switch char {
		case '"':
			inString = true // 进入字符串模式
		case '{', '[':
			stack = append(stack, char) // 入栈
		case '}', ']':
			if len(stack) == 0 {
				return false // 栈空，右括号多了
			}
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1] // Pop

			// 检查匹配
			if (char == '}' && top != '{') || (char == ']' && top != '[') {
				return false
			}
		}
		// 其他字符（: , 0-9 a-z）在结构验证中忽略
	}

	// 最终检查：栈必须为空，且不能停留在字符串模式中
	return len(stack) == 0 && !inString
}

func main() {
	validCases := []string{
		`{"name": "Alice", "age": 30, "scores": [100, 99]}`,
		`{"key": "value with } bracket inside"}`, // 字符串内的 } 不应干扰
		`{"key": "quote \" in string"}`,          // 转义引号
		`[{}, {"a": [1, 2]}]`,
	}

	invalidCases := []string{
		`{"name": "Alice"`,        // 缺右括号
		`{"key": [1, 2}}`,         // 括号不匹配
		`{"key": "unfinished str`, // 字符串未闭合
		`[1, 2]]`,                 // 多余右括号
	}

	fmt.Println("--- Valid Cases ---")
	for _, s := range validCases {
		if ValidateJSONStructure(s) {
			fmt.Println("PASS:", s)
		} else {
			fmt.Println("FAIL (Unexpected):", s)
		}
	}

	fmt.Println("\n--- Invalid Cases ---")
	for _, s := range invalidCases {
		if !ValidateJSONStructure(s) {
			fmt.Println("PASS (Correctly Rejected):", s)
		} else {
			fmt.Println("FAIL (Should be invalid):", s)
		}
	}
}
