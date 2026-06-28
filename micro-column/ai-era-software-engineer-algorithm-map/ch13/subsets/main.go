package main

import (
	"fmt"
	"math"
)

func subsets(nums []int) [][]int {
	n := len(nums)
	// 总子集数 = 2^n
	totalSubsets := int(math.Pow(2, float64(n)))
	res := make([][]int, 0, totalSubsets)

	// 遍历 0 到 2^n - 1 的每个 mask
	for mask := 0; mask < totalSubsets; mask++ {
		var subset []int
		// 检查 mask 的每一位
		for i := 0; i < n; i++ {
			// 如果 mask 的第 i 位是 1
			if (mask>>i)&1 == 1 {
				subset = append(subset, nums[i])
			}
		}
		res = append(res, subset)
	}
	return res
}

func main() {
	nums := []int{1, 2, 3}
	fmt.Printf("Subsets of %v:\n%v\n", nums, subsets(nums))
}
