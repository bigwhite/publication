package main

import "fmt"

func nearestExit(maze [][]byte, entrance []int) int {
	rows, cols := len(maze), len(maze[0])
	// 方向数组：上右下左
	dirs := [][]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

	// 队列：存坐标 [row, col, steps]
	queue := [][]int{{entrance[0], entrance[1], 0}}

	// 标记起点已访问，避免走回头路
	maze[entrance[0]][entrance[1]] = '+'

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		r, c, steps := curr[0], curr[1], curr[2]

		// 遍历四个方向
		for _, d := range dirs {
			nr, nc := r+d[0], c+d[1]

			// 检查边界和墙
			if nr >= 0 && nr < rows && nc >= 0 && nc < cols && maze[nr][nc] == '.' {
				// 核心判断：是否到达边界且不是入口？
				if nr == 0 || nr == rows-1 || nc == 0 || nc == cols-1 {
					return steps + 1
				}

				// 标记已访问并入队
				maze[nr][nc] = '+'
				queue = append(queue, []int{nr, nc, steps + 1})
			}
		}
	}

	return -1
}

func main() {
	maze := [][]byte{
		{'+', '+', '.', '+'},
		{'.', '.', '.', '+'},
		{'+', '+', '+', '.'},
	}
	ent := []int{1, 2}
	fmt.Printf("Steps to Exit: %d\n", nearestExit(maze, ent))
}
