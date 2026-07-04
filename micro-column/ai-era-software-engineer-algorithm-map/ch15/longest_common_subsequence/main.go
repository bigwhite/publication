package main

import "fmt"

func longestCommonSubsequence(text1 string, text2 string) int {
	m, n := len(text1), len(text2)
	// dp[i][j] 对应 text1 的前 i 个字符和 text2 的前 j 个字符
	// 尺寸 +1 是为了处理空串的情况 (padding)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			// 注意：字符串索引从 0 开始，所以对应的是 i-1 和 j-1
			if text1[i-1] == text2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}

	return dp[m][n]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	t1 := "abcde"
	t2 := "ace"
	fmt.Printf("LCS('%s', '%s') = %d\n", t1, t2, longestCommonSubsequence(t1, t2))
}
