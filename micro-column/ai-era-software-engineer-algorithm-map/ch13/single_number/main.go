package main

import "fmt"

func singleNumber(nums []int) int {
	res := 0
	for _, n := range nums {
		// 异或操作
		// 0 ^ n = n
		// n ^ n = 0
		// a ^ b ^ a = b
		res ^= n
	}
	return res
}

func main() {
	nums := []int{4, 1, 2, 1, 2}
	fmt.Printf("Nums: %v, Single: %d\n", nums, singleNumber(nums))
}
