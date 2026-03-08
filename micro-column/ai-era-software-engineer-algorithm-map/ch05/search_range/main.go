package main

import "fmt"

// searchRange 主函数
func searchRange(nums []int, target int) []int {
	left := findFirst(nums, target)
	// 如果左边界都没找到，说明数组里根本没有 target
	if left == -1 {
		return []int{-1, -1}
	}
	right := findLast(nums, target)
	return []int{left, right}
}

// 寻找第一个等于 target 的位置 (Lower Bound)
func findFirst(nums []int, target int) int {
	left, right := 0, len(nums)-1
	res := -1
	for left <= right {
		mid := left + (right-left)/2
		if nums[mid] == target {
			res = mid
			// 核心：找到了不要停，继续向左尝试，看前面还有没有
			right = mid - 1
		} else if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return res
}

// 寻找最后一个等于 target 的位置
func findLast(nums []int, target int) int {
	left, right := 0, len(nums)-1
	res := -1
	for left <= right {
		mid := left + (right-left)/2
		if nums[mid] == target {
			res = mid
			// 核心：找到了不要停，继续向右尝试
			left = mid + 1
		} else if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return res
}

func main() {
	nums := []int{5, 7, 7, 8, 8, 10}
	target := 8
	fmt.Printf("Nums: %v, Target: %d, Range: %v\n", nums, target, searchRange(nums, target))
}
