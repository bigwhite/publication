package main

import "fmt"

func rotate(matrix [][]int) {
	n := len(matrix)

	// 1. 转置 (Transpose)
	// 只需要遍历对角线右上方的元素
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ { // 注意 j 从 i+1 开始
			matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
		}
	}

	// 2. 左右翻转 (Reverse each row)
	for i := 0; i < n; i++ {
		left, right := 0, n-1
		for left < right {
			matrix[i][left], matrix[i][right] = matrix[i][right], matrix[i][left]
			left++
			right--
		}
	}
}

func printMatrix(m [][]int) {
	for _, row := range m {
		fmt.Println(row)
	}
}

func main() {
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	fmt.Println("Original:")
	printMatrix(matrix)

	rotate(matrix)

	fmt.Println("Rotated 90 degrees:")
	printMatrix(matrix)
}
