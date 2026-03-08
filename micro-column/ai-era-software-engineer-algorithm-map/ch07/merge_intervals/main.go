package main

import (
	"fmt"
	"math"
	"sort"
)

func merge(intervals [][]int) [][]int {
	if len(intervals) <= 1 {
		return intervals
	}

	// 1. 贪心第一步：按 start 排序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	var res [][]int
	// 先把第一个放入结果集
	res = append(res, intervals[0])

	for i := 1; i < len(intervals); i++ {
		curr := intervals[i]
		last := res[len(res)-1] // 引用结果集中的最后一个区间

		// 2. 贪心决策：如果当前区间起点 <= 上一个区间终点 -> 重叠
		if curr[0] <= last[1] {
			// 合并：更新上一个区间的终点为两者的最大值
			// 注意：必须直接修改 res 中的元素
			res[len(res)-1][1] = int(math.Max(float64(last[1]), float64(curr[1])))
		} else {
			// 不重叠，直接加入
			res = append(res, curr)
		}
	}

	return res
}

func main() {
	intervals := [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	fmt.Printf("Original: %v\n", intervals)
	fmt.Printf("Merged:   %v\n", merge(intervals))
}
