package main

import "fmt"

// isMatch 实现简单的正则匹配 (. 和 *)
func isMatch(text string, pattern string) bool {
	// 1. 终止条件：pattern 耗尽
	if len(pattern) == 0 {
		return len(text) == 0
	}

	// 检查首字符是否匹配（注意 text 可能为空）
	firstMatch := len(text) > 0 && (pattern[0] == text[0] || pattern[0] == '.')

	// 2. 处理 '*' (Lookahead)
	if len(pattern) >= 2 && pattern[1] == '*' {
		// 两种选择（回溯的分叉点）：
		// Option A: '*' 匹配 0 次，跳过 pattern 前两个字符 (x*)
		// Option B: '*' 匹配 1+ 次（前提是 firstMatch），text 前进 1，pattern 保持不变
		return isMatch(text, pattern[2:]) || (firstMatch && isMatch(text[1:], pattern))
	} else {
		// 3. 普通匹配
		return firstMatch && isMatch(text[1:], pattern[1:])
	}
}

func main() {
	cases := []struct {
		text    string
		pattern string
		want    bool
	}{
		{"aa", "a", false},
		{"aa", "a*", true},
		{"ab", ".*", true},
		{"aab", "c*a*b", true},
		{"mississippi", "mis*is*p*.", false},
	}

	fmt.Println("--- Simple Regex Engine ---")
	for _, c := range cases {
		got := isMatch(c.text, c.pattern)
		status := "FAIL"
		if got == c.want {
			status = "PASS"
		}
		fmt.Printf("[%s] Text: %-12s Pat: %-10s -> %v\n", status, c.text, c.pattern, got)
	}
}
