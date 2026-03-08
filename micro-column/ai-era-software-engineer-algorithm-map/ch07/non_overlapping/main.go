package main

import (
	"fmt"
	"sort"
)

func eraseOverlapIntervals(intervals [][]int) int {
	if len(intervals) == 0 {
		return 0
	}

	// 1. 贪心核心：按 End 升序排序
	// 只有结束得早，才不容易挡后面人的路
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][1] < intervals[j][1]
	})

	// count 记录最多能保留的不重叠区间数
	count := 1
	// 记录上一个被选中区间的结束时间
	lastEnd := intervals[0][1]

	for i := 1; i < len(intervals); i++ {
		// 2. 决策：如果当前区间开始时间 >= 上一个结束时间
		if intervals[i][0] >= lastEnd {
			// 这是一个兼容的好区间，选中它
			count++
			lastEnd = intervals[i][1]
		}
		// 否则：产生冲突了。
		// 贪心逻辑：既然 intervals[i] 结束得比 lastEnd 晚（因为排序了），
		// 或者即便结束得早但已经和 lastEnd 冲突了，我们优先保留 lastEnd (因为它已经是最优选)
		// 所以这里实际上是丢弃了 intervals[i]
	}

	// 需要移除的数量 = 总数 - 最多保留数
	return len(intervals) - count
}

func main() {
	intervals := [][]int{{1, 2}, {2, 3}, {3, 4}, {1, 3}}
	fmt.Printf("Intervals: %v\n", intervals)
	fmt.Printf("Min Removed: %d\n", eraseOverlapIntervals(intervals))
}
