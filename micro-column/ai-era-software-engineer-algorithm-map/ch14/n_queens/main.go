package main

import (
	"fmt"
	"strings"
)

func solveNQueens(n int) [][]string {
	var res [][]string
	// board 记录每一行皇后所在的列索引。board[row] = col
	board := make([]int, n)
	for i := range board {
		board[i] = -1
	}

	// 辅助检查函数
	isValid := func(row, col int) bool {
		for r := 0; r < row; r++ {
			c := board[r]
			// 1. 同列检查
			if c == col {
				return false
			}
			// 2. 斜线检查：行差 == 列差
			if r-c == row-col || r+c == row+col {
				return false
			}
		}
		return true
	}

	var backtrack func(row int)
	backtrack = func(row int) {
		// 1. 终止条件：放完了最后一行
		if row == n {
			res = append(res, generateBoard(board, n))
			return
		}

		// 2. 遍历当前行的每一列
		for col := 0; col < n; col++ {
			// 剪枝
			if !isValid(row, col) {
				continue
			}

			// 3. 做选择
			board[row] = col

			// 4. 递归下一行
			backtrack(row + 1)

			// 5. 撤销选择
			board[row] = -1
		}
	}

	backtrack(0)
	return res
}

// 辅助函数：将 board 数组转为题目要求的字符串格式
func generateBoard(board []int, n int) []string {
	var s []string
	for _, col := range board {
		rowStr := strings.Repeat(".", n)
		rowStr = rowStr[:col] + "Q" + rowStr[col+1:]
		s = append(s, rowStr)
	}
	return s
}

func main() {
	n := 4
	solutions := solveNQueens(n)
	fmt.Printf("N=%d, Solutions count: %d\n", n, len(solutions))
	for i, sol := range solutions {
		fmt.Printf("Solution %d:\n", i+1)
		for _, row := range sol {
			fmt.Println(row)
		}
		fmt.Println()
	}
}
