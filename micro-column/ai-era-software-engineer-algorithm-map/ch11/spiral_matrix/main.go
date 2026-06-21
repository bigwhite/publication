package main

import "fmt"

func spiralOrder(matrix [][]int) []int {
	if len(matrix) == 0 {
		return nil
	}

	m, n := len(matrix), len(matrix[0])
	res := make([]int, 0, m*n)

	top, bottom := 0, m-1
	left, right := 0, n-1

	for top <= bottom && left <= right {
		// 1. 向右 (top row)
		for j := left; j <= right; j++ {
			res = append(res, matrix[top][j])
		}
		top++ // 上边界收缩

		// 2. 向下 (right col)
		for i := top; i <= bottom; i++ {
			res = append(res, matrix[i][right])
		}
		right-- // 右边界收缩

		// 检查是否越界（防止单行/单列矩阵重复遍历）
		if top <= bottom {
			// 3. 向左 (bottom row)
			for j := right; j >= left; j-- {
				res = append(res, matrix[bottom][j])
			}
			bottom-- // 下边界收缩
		}

		if left <= right {
			// 4. 向上 (left col)
			for i := bottom; i >= top; i-- {
				res = append(res, matrix[i][left])
			}
			left++ // 左边界收缩
		}
	}
	return res
}

func main() {
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	fmt.Println("Spiral:", spiralOrder(matrix))
}
