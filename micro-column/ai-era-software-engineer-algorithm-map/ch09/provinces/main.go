package main

import "fmt"

// UnionFind 并查集结构
type UnionFind struct {
	parent []int
	count  int // 连通分量个数
}

func NewUnionFind(n int) *UnionFind {
	parent := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i // 初始时，每个人的老大是自己
	}
	return &UnionFind{parent: parent, count: n}
}

// Find 查找老大（带路径压缩）
func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] != x {
		// 递归查找，并顺手把沿途节点的 parent 直接指向老大
		uf.parent[x] = uf.Find(uf.parent[x])
	}
	return uf.parent[x]
}

// Union 合并两个集合
func (uf *UnionFind) Union(x, y int) {
	rootX := uf.Find(x)
	rootY := uf.Find(y)
	if rootX != rootY {
		uf.parent[rootX] = rootY // X 归顺 Y
		uf.count--               // 集合少了一个
	}
}

func findCircleNum(isConnected [][]int) int {
	n := len(isConnected)
	uf := NewUnionFind(n)

	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if isConnected[i][j] == 1 {
				uf.Union(i, j)
			}
		}
	}
	return uf.count
}

func main() {
	// 1-2 相连, 3 独立
	matrix := [][]int{
		{1, 1, 0},
		{1, 1, 0},
		{0, 0, 1},
	}
	fmt.Printf("Provinces: %d\n", findCircleNum(matrix)) // Output: 2
}
