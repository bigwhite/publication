package main

import "fmt"

func permute(nums []int) [][]int {
	var res [][]int
	var path []int

	// used 记录数字是否已在当前 path 中
	used := make([]bool, len(nums))

	var backtrack func()
	backtrack = func() {
		// 1. 终止条件：path 长度等于 nums，说明凑齐了一个排列
		if len(path) == len(nums) {
			// 注意：Go 的切片是引用，必须 copy 一份
			temp := make([]int, len(path))
			copy(temp, path)
			res = append(res, temp)
			return
		}

		// 2. 遍历选择
		for i, num := range nums {
			// 剪枝：如果已经选过，跳过
			if used[i] {
				continue
			}

			// 3. 做选择
			path = append(path, num)
			used[i] = true

			// 4. 递归
			backtrack()

			// 5. 撤销选择 (回溯)
			used[i] = false
			path = path[:len(path)-1]
		}
	}

	backtrack()
	return res
}

func main() {
	nums := []int{1, 2, 3}
	fmt.Printf("Permutations: %v\n", permute(nums))
}
