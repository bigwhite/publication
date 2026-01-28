package main

import "fmt"

func TwoSum(numbers []int, target int) []int {
	left, right := 0, len(numbers)-1

	for left < right {
		sum := numbers[left] + numbers[right]
		if sum == target {
			// 题目要求下标从 1 开始
			return []int{left + 1, right + 1}
		} else if sum < target {
			// 和太小了，需要变大。
			// 因为数组有序，只有右移左指针才能让和变大
			left++
		} else {
			// 和太大了，需要变小。
			// 左移右指针
			right--
		}
	}
	return nil
}

func main() {
	nums := []int{2, 7, 11, 15}
	target := 9
	result := TwoSum(nums, target)
	fmt.Println("Indices:", result) // Output: [1, 2]
}
