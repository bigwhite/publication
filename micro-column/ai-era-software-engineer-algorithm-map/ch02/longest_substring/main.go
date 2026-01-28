package main

import "fmt"

func LengthOfLongestSubstring(s string) int {
	// window 记录字符上一次出现的索引
	// 使用 [256]int 替代 map，对于 ASCII 字符集性能更好
	window := [256]int{}
	for i := range window {
		window[i] = -1 // 初始化为 -1
	}

	left, maxLen := 0, 0

	for right := 0; right < len(s); right++ {
		charIndex := s[right]

		// 如果当前字符之前出现过，并且出现在当前窗口内 (>= left)
		// 说明窗口内有重复，left 需要跳跃
		if prevIdx := window[charIndex]; prevIdx >= left {
			left = prevIdx + 1
		}

		// 更新当前字符的最新位置
		window[charIndex] = right

		// 更新最大长度
		if currLen := right - left + 1; currLen > maxLen {
			maxLen = currLen
		}
	}
	return maxLen
}

func main() {
	s := "abcabcbb"
	fmt.Printf("String: %s, Max Length: %d\n", s, LengthOfLongestSubstring(s))
}
