package main

import (
	"fmt"
	"math"
)

func MinSubArrayLen(target int, nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}

	left := 0
	sum := 0
	minLen := math.MaxInt32

	for right := 0; right < n; right++ {
		sum += nums[right] // 1. 进窗：累加

		// 2. 满足条件时，尝试收缩窗口
		for sum >= target {
			currLen := right - left + 1
			if currLen < minLen {
				minLen = currLen
			}

			// 3. 出窗：移除左边元素
			sum -= nums[left]
			left++
		}
	}

	if minLen == math.MaxInt32 {
		return 0
	}
	return minLen
}

func main() {
	target := 7
	nums := []int{2, 3, 1, 2, 4, 3}
	fmt.Printf("Target: %d, Nums: %v, Min Len: %d\n", target, nums, MinSubArrayLen(target, nums))
}
