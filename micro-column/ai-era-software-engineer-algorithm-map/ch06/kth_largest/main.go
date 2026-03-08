package main

import (
	"container/heap"
	"fmt"
)

// IntMinHeap 定义一个小顶堆
type IntMinHeap []int

func (h IntMinHeap) Len() int           { return len(h) }
func (h IntMinHeap) Less(i, j int) bool { return h[i] < h[j] } // 小顶堆
func (h IntMinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntMinHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *IntMinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func FindKthLargest(nums []int, k int) int {
	h := &IntMinHeap{}
	heap.Init(h)

	for _, num := range nums {
		// 1. 先把前 k 个填满
		if h.Len() < k {
			heap.Push(h, num)
		} else if num > (*h)[0] {
			// 2. 如果新元素比堆顶（这 k 个里最小的）大
			// 说明新元素更有资格成为 Top K
			// 替换堆顶：先 Pop 最小的，再 Push 新的
			heap.Pop(h)
			heap.Push(h, num)

			// 优化技巧：
			// 其实可以直接 (*h)[0] = num 然后 heap.Fix(h, 0)
			// 这样少一次 Slice 的 resize 操作，性能更好
		}
	}

	// 堆顶就是第 K 大
	return (*h)[0]
}

func main() {
	nums := []int{3, 2, 1, 5, 6, 4}
	k := 2
	fmt.Printf("Nums: %v, %d-th Largest: %d\n", nums, k, FindKthLargest(nums, k))
}
