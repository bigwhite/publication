package main

import (
	"fmt"
)

func coinChange(coins []int, amount int) int {
	// dp[i] 表示凑成金额 i 所需的最少硬币数
	// 初始化为一个极大值 (amount + 1)，方便求 min
	dp := make([]int, amount+1)
	for i := range dp {
		dp[i] = amount + 1
	}

	// Base Case: 凑成金额 0 需要 0 个硬币
	dp[0] = 0

	// 遍历每一个金额状态
	for i := 1; i <= amount; i++ {
		// 尝试每一枚硬币
		for _, coin := range coins {
			if i-coin >= 0 {
				// 状态转移：当前最少 = min(当前记录, 凑成(i-coin)的最少 + 1枚当前硬币)
				if dp[i-coin]+1 < dp[i] {
					dp[i] = dp[i-coin] + 1
				}
			}
		}
	}

	if dp[amount] > amount {
		return -1
	}
	return dp[amount]
}

func main() {
	coins := []int{1, 2, 5}
	amount := 11
	fmt.Printf("Coins: %v, Amount: %d, Min Coins: %d\n", coins, amount, coinChange(coins, amount))
}
