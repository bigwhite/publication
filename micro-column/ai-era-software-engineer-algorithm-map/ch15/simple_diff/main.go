package main

import (
	"fmt"
	"strings"
)

// DiffLine 定义差异行
type DiffLine struct {
	Type    string // " ", "+", "-"
	Content string
}

// ComputeDiff 计算两个文本的差异
func ComputeDiff(oldLines, newLines []string) []DiffLine {
	m, n := len(oldLines), len(newLines)

	// 1. 构建 DP 表 (LCS)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if oldLines[i-1] == newLines[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				if dp[i-1][j] >= dp[i][j-1] {
					dp[i][j] = dp[i-1][j]
				} else {
					dp[i][j] = dp[i][j-1]
				}
			}
		}
	}

	// 2. 回溯路径，生成 Diff
	var diffs []DiffLine
	i, j := m, n
	for i > 0 || j > 0 {
		if i > 0 && j > 0 && oldLines[i-1] == newLines[j-1] {
			// Case 1: 相同行 (Unchanged)
			diffs = append(diffs, DiffLine{" ", oldLines[i-1]})
			i--
			j--
		} else if j > 0 && (i == 0 || dp[i][j-1] >= dp[i-1][j]) {
			// Case 2: 优先看向左边 (dp[i][j-1]) -> NewLines 有新内容 -> Add (+)
			// 注意：这里回溯条件的判断需要细心，通常优先处理 Add 或 Del 都可以
			diffs = append(diffs, DiffLine{"+", newLines[j-1]})
			j--
		} else if i > 0 && (j == 0 || dp[i][j-1] < dp[i-1][j]) {
			// Case 3: 向上看 (dp[i-1][j]) -> OldLines 有内容没匹配 -> Del (-)
			diffs = append(diffs, DiffLine{"-", oldLines[i-1]})
			i--
		}
	}

	// 因为是逆向回溯，需要反转结果
	reverse(diffs)
	return diffs
}

func reverse(s []DiffLine) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func main() {
	textA := `
func main() {
	fmt.Println("Hello")
	fmt.Println("World")
}
`
	textB := `
func main() {
	fmt.Println("Hello")
	fmt.Println("Go")
}
`
	// 简单的按行分割
	linesA := strings.Split(strings.TrimSpace(textA), "\n")
	linesB := strings.Split(strings.TrimSpace(textB), "\n")

	diffs := ComputeDiff(linesA, linesB)

	fmt.Println("--- Diff Result ---")
	for _, d := range diffs {
		fmt.Printf("%s %s\n", d.Type, d.Content)
	}
}
