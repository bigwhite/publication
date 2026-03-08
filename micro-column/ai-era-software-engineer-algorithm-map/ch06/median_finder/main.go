package main

import (
	"container/heap"
	"fmt"
)

// --- Min Heap ---
type MinHeap []int

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// --- Max Heap (Go 默认只有小顶堆，取反实现大顶堆) ---
type MaxHeap []int

func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i] > h[j] } // 核心区别：>
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type MedianFinder struct {
	low  *MaxHeap // 存较小的一半
	high *MinHeap // 存较大的一半
}

func Constructor() MedianFinder {
	return MedianFinder{
		low:  &MaxHeap{},
		high: &MinHeap{},
	}
}

func (this *MedianFinder) AddNum(num int) {
	// 1. 先放入 low，然后把 low 最大的给 high，保证顺序
	heap.Push(this.low, num)
	heap.Push(this.high, heap.Pop(this.low))

	// 2. 平衡大小：保证 low 的数量 >= high 的数量
	if this.low.Len() < this.high.Len() {
		heap.Push(this.low, heap.Pop(this.high))
	}
}

func (this *MedianFinder) FindMedian() float64 {
	if this.low.Len() > this.high.Len() {
		return float64((*this.low)[0])
	}
	return float64((*this.low)[0]+(*this.high)[0]) / 2.0
}

func main() {
	mf := Constructor()
	mf.AddNum(1)
	mf.AddNum(2)
	fmt.Println("Median (1,2):", mf.FindMedian()) // 1.5
	mf.AddNum(3)
	fmt.Println("Median (1,2,3):", mf.FindMedian()) // 2
}
