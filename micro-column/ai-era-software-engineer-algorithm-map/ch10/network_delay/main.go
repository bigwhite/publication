package main

import (
	"container/heap"
	"fmt"
	"math"
)

// 优先队列项
type Item struct {
	Node int
	Dist int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].Dist < pq[j].Dist } // Min Heap
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*Item)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func networkDelayTime(times [][]int, n int, k int) int {
	// 1. 建图：邻接表
	graph := make(map[int][][]int)
	for _, t := range times {
		u, v, w := t[0], t[1], t[2]
		graph[u] = append(graph[u], []int{v, w})
	}

	// 2. 初始化距离表
	dist := make(map[int]int)
	for i := 1; i <= n; i++ {
		dist[i] = math.MaxInt32
	}
	dist[k] = 0

	// 3. Dijkstra 主循环
	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, &Item{Node: k, Dist: 0})

	for pq.Len() > 0 {
		curr := heap.Pop(pq).(*Item)
		u, d := curr.Node, curr.Dist

		// 剪枝：如果当前距离已经比记录的短路径长，跳过（Lazy Deletion）
		if d > dist[u] {
			continue
		}

		// 松弛邻居
		for _, edge := range graph[u] {
			v, w := edge[0], edge[1]
			if dist[u]+w < dist[v] {
				dist[v] = dist[u] + w
				heap.Push(pq, &Item{Node: v, Dist: dist[v]})
			}
		}
	}

	// 4. 统计结果
	maxDist := 0
	for i := 1; i <= n; i++ {
		if dist[i] == math.MaxInt32 {
			return -1 // 有节点不可达
		}
		if dist[i] > maxDist {
			maxDist = dist[i]
		}
	}
	return maxDist
}

func main() {
	// 2->1(1), 2->3(1), 3->4(1)
	times := [][]int{{2, 1, 1}, {2, 3, 1}, {3, 4, 1}}
	n, k := 4, 2
	fmt.Printf("Max Delay: %d\n", networkDelayTime(times, n, k))
}
