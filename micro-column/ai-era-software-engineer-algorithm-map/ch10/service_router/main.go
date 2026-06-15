package main

import (
	"container/heap"
	"fmt"
	"math"
	"strings"
)

// Edge 定义服务调用边
type Edge struct {
	To      string
	Latency int // 毫秒
}

// Graph 服务依赖图
type Graph map[string][]Edge

// PathResult 路径结果
type PathResult struct {
	Path         []string
	TotalLatency int
}

// --- Priority Queue ---
type Item struct {
	Node    string
	Latency int
	index   int
}
type PQ []*Item

func (pq PQ) Len() int           { return len(pq) }
func (pq PQ) Less(i, j int) bool { return pq[i].Latency < pq[j].Latency }
func (pq PQ) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i]; pq[i].index = i; pq[j].index = j }
func (pq *PQ) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// FindShortestPath 寻找最低延迟路径
func FindShortestPath(g Graph, start, end string) *PathResult {
	// distMap: 到达某节点的最短距离
	distMap := make(map[string]int)
	// prevMap: 记录前驱节点，用于回溯路径
	prevMap := make(map[string]string)

	for node := range g {
		distMap[node] = math.MaxInt32
	}
	distMap[start] = 0

	pq := &PQ{}
	heap.Init(pq)
	heap.Push(pq, &Item{Node: start, Latency: 0})

	found := false

	for pq.Len() > 0 {
		curr := heap.Pop(pq).(*Item)
		u := curr.Node

		if u == end {
			found = true
			break
		}

		if curr.Latency > distMap[u] {
			continue
		}

		for _, edge := range g[u] {
			v := edge.To
			newDist := distMap[u] + edge.Latency

			// 松弛操作
			if newDist < distMap[v] {
				distMap[v] = newDist
				prevMap[v] = u
				heap.Push(pq, &Item{Node: v, Latency: newDist})
			}
		}
	}

	if !found {
		return nil
	}

	// 回溯路径
	path := []string{}
	curr := end
	for curr != "" {
		path = append([]string{curr}, path...) // Prepend
		if curr == start {
			break
		}
		curr = prevMap[curr]
	}

	return &PathResult{
		Path:         path,
		TotalLatency: distMap[end],
	}
}

func main() {
	// 构建服务网络
	// Gateway -> Auth (10ms)
	// Gateway -> Cache (5ms)
	// Auth -> Core (20ms)
	// Cache -> Core (10ms)
	// Core -> DB (50ms)
	// Cache -> DB (80ms) -- 假设直连DB很慢

	graph := Graph{
		"Gateway": {{To: "Auth", Latency: 10}, {To: "Cache", Latency: 5}},
		"Auth":    {{To: "Core", Latency: 20}},
		"Cache":   {{To: "Core", Latency: 10}, {To: "DB", Latency: 80}},
		"Core":    {{To: "DB", Latency: 50}},
		"DB":      {},
	}

	start, end := "Gateway", "DB"
	result := FindShortestPath(graph, start, end)

	if result != nil {
		fmt.Printf("Optimal Route: %s\n", strings.Join(result.Path, " -> "))
		fmt.Printf("Total Latency: %dms\n", result.TotalLatency)
	} else {
		fmt.Println("No route found!")
	}
}
