package main

import "fmt"

// RemoveDuplicates 使用快慢指针原地去重
// 输入: sorted slice
// 输出: 新的长度
func RemoveDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	// slow 指针：指向当前“无重复序列”的最后一个位置
	slow := 0

	// fast 指针：在前探路，寻找与 slow 位置不同的新元素
	for fast := 1; fast < len(nums); fast++ {
		if nums[fast] != nums[slow] {
			// 发现了新大陆（新元素）
			slow++
			// 将新元素挪到 slow 的位置
			// 这一步是关键：我们在原地复写，没有内存分配
			nums[slow] = nums[fast]
		}
		// 如果 nums[fast] == nums[slow]，说明是重复元素，
		// fast 继续往前跑，slow 原地不动，相当于跳过了重复项
	}

	// 长度是索引 + 1
	return slow + 1
}

func main() {
	data := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	newLen := RemoveDuplicates(data)
	// 注意：我们只关心前 newLen 个元素
	fmt.Printf("New Length: %d, Data: %v\n", newLen, data[:newLen])
}
