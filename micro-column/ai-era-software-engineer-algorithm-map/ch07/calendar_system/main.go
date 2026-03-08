package main

import (
	"fmt"
)

// Booking 代表一次预订
type Booking struct {
	Start, End int
}

// Calendar 日历系统
type Calendar struct {
	bookings []Booking // 始终保持按 Start 排序
}

func NewCalendar() *Calendar {
	return &Calendar{
		bookings: make([]Booking, 0),
	}
}

// Book 尝试预订，返回是否成功
// 时间复杂度：O(N) - 主要是切片插入的开销；查找本身是 O(log N)
// 如果使用真正的平衡二叉搜索树，整体复杂度可降为 O(log N)
func (c *Calendar) Book(start, end int) bool {
	// 1. 二分查找：找到第一个 Start >= 新 Start 的预订位置
	// 我们需要检查该位置的前一个和后一个，看是否冲突

	// 这里为了代码简洁，且 Go 没有内置 upper_bound，我们手动实现二分逻辑
	// 目标：找到 idx，使得 bookings[idx].Start >= start
	left, right := 0, len(c.bookings)
	for left < right {
		mid := left + (right-left)/2
		if c.bookings[mid].Start >= start {
			right = mid
		} else {
			left = mid + 1
		}
	}
	idx := left // 插入位置

	// 2. 冲突检测 (贪心检查：只看相邻的)

	// 检查前一个预订 (idx-1)：如果前一个的 End > 现在的 Start -> 冲突
	if idx > 0 {
		prevBooking := c.bookings[idx-1]
		if prevBooking.End > start {
			return false
		}
	}

	// 检查后一个预订 (idx)：如果现在的 End > 后一个的 Start -> 冲突
	if idx < len(c.bookings) {
		nextBooking := c.bookings[idx]
		if end > nextBooking.Start {
			return false
		}
	}

	// 3. 无冲突，插入并保持有序
	// Go 切片插入模板代码
	c.bookings = append(c.bookings, Booking{})        // 扩容
	copy(c.bookings[idx+1:], c.bookings[idx:])        // 后移
	c.bookings[idx] = Booking{Start: start, End: end} // 填入

	return true
}

func main() {
	cal := NewCalendar()

	events := [][]int{
		{10, 20}, // Book 1: OK
		{15, 25}, // Book 2: Conflict with [10, 20)
		{20, 30}, // Book 3: OK (紧接 [10, 20))
		{5, 10},  // Book 4: OK (在 [10, 20) 之前)
		{25, 35}, // Book 5: Conflict with [20, 30)
	}

	fmt.Println("--- Booking Log ---")
	for _, e := range events {
		success := cal.Book(e[0], e[1])
		status := "Accepted"
		if !success {
			status = "Rejected"
		}
		fmt.Printf("Book [%d, %d): %s\n", e[0], e[1], status)
	}

	fmt.Println("\n--- Final Schedule ---")
	for _, b := range cal.bookings {
		fmt.Printf("[%d, %d) ", b.Start, b.End)
	}
	fmt.Println()
}
