package tree

import (
	"container/heap"
)

// PriorityQueue 是一个由 Node 组成的最小堆
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

// Less 决定排序规则：频率小的排在前面
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Freq < pq[j].Freq
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	node := x.(*Node)
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	*pq = old[0 : n-1]
	return node
}

// NewQueue 辅助函数：初始化并构建堆
func NewQueue(nodes []*Node) *PriorityQueue {
	pq := PriorityQueue(nodes)
	heap.Init(&pq)
	return &pq
}
