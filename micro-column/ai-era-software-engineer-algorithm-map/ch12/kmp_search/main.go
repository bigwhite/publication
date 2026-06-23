package main

import "fmt"

// buildNext 构建 Next 数组 (PMT: Partial Match Table)
// next[i] 表示 pattern[0...i] 的最长公共前后缀长度
func buildNext(pattern string) []int {
	m := len(pattern)
	next := make([]int, m)
	j := 0 // j 是前缀的末尾，也是当前最长公共前后缀的长度

	// i 从 1 开始遍历后缀
	for i := 1; i < m; i++ {
		// 1. 不匹配，回退 j (核心逻辑)
		for j > 0 && pattern[i] != pattern[j] {
			j = next[j-1] // 回退到上一级最长公共前后缀的位置
		}
		// 2. 匹配，j 前进
		if pattern[i] == pattern[j] {
			j++
		}
		// 3. 记录当前位置的 next 值
		next[i] = j
	}
	return next
}

// KMP 搜索算法
func StrStrKMP(haystack, needle string) int {
	if len(needle) == 0 {
		return 0
	}

	next := buildNext(needle)
	j := 0 // needle 的指针

	for i := 0; i < len(haystack); i++ {
		// 1. 不匹配，根据 next 数组回退 j
		// 注意：这里是 while 循环，可能多次回退
		for j > 0 && haystack[i] != needle[j] {
			j = next[j-1]
		}
		// 2. 匹配，j 前进
		if haystack[i] == needle[j] {
			j++
		}
		// 3. 完全匹配
		if j == len(needle) {
			return i - len(needle) + 1
		}
	}

	return -1
}

func main() {
	text := "ABABABCA"
	pat := "ABABC"
	idx := StrStrKMP(text, pat)
	fmt.Printf("Text: %s, Pattern: %s, Index: %d\n", text, pat, idx)
}
