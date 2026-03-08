package main

import (
	"fmt"
)

func minEatingSpeed(piles []int, h int) int {
	// 1. 确定二分范围 [left, right]
	// 最小速度是 1，最大速度是香蕉堆里的最大值（因为一小时最多吃一堆）
	maxPile := 0
	for _, p := range piles {
		if p > maxPile {
			maxPile = p
		}
	}

	left, right := 1, maxPile
	ans := maxPile

	for left <= right {
		mid := left + (right-left)/2

		// 2. Check 函数：以速度 mid 能否在 h 小时内吃完？
		if canFinish(piles, h, mid) {
			// 能吃完，记录答案，并尝试更小的速度（向左收缩）
			ans = mid
			right = mid - 1
		} else {
			// 吃不完，说明速度太慢了，需要加速（向右收缩）
			left = mid + 1
		}
	}
	return ans
}

// 辅助函数：计算以速度 k 吃完所有香蕉需要的小时数
func canFinish(piles []int, h int, k int) bool {
	hoursUsed := 0
	for _, p := range piles {
		// 向上取整：p/k
		// 技巧：(p + k - 1) / k 等价于 math.Ceil(float(p)/float(k))
		hoursUsed += (p + k - 1) / k
	}
	return hoursUsed <= h
}

func main() {
	piles := []int{3, 6, 7, 11}
	h := 8
	fmt.Printf("Piles: %v, Hours: %d, Min Speed: %d\n", piles, h, minEatingSpeed(piles, h))
}
