package main

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
)

// Ring is the consistent hash ring
type Ring struct {
	nodes        map[string]bool // 物理节点
	virtualCount int             // 每个物理节点的虚拟节点数
	keys         []int           // 排序后的虚拟节点哈希值 (0-359)
	vnodeMap     map[int]string  // 虚拟节点哈希 -> 物理节点名称
}

func NewRing(virtualCount int) *Ring {
	return &Ring{
		nodes:        make(map[string]bool),
		virtualCount: virtualCount,
		keys:         []int{},
		vnodeMap:     make(map[int]string),
	}
}

// hash a string to an angle (0-359)
func (r *Ring) hash(key string) int {
	return int(crc32.ChecksumIEEE([]byte(key))) % 360
}

// AddNode adds a physical node to the ring
func (r *Ring) AddNode(node string) {
	if _, ok := r.nodes[node]; ok {
		return
	}
	r.nodes[node] = true
	for i := 0; i < r.virtualCount; i++ {
		vnodeKey := node + "#" + strconv.Itoa(i)
		hash := r.hash(vnodeKey)
		r.keys = append(r.keys, hash)
		r.vnodeMap[hash] = node
	}
	sort.Ints(r.keys)
}

// RemoveNode removes a physical node from the ring
func (r *Ring) RemoveNode(node string) {
	if _, ok := r.nodes[node]; !ok {
		return
	}
	delete(r.nodes, node)
	for i := 0; i < r.virtualCount; i++ {
		vnodeKey := node + "#" + strconv.Itoa(i)
		hash := r.hash(vnodeKey)
		delete(r.vnodeMap, hash)
	}
	r.keys = []int{}
	for k := range r.vnodeMap {
		r.keys = append(r.keys, k)
	}
	sort.Ints(r.keys)
}

// GetNode returns the physical node for a given key
func (r *Ring) GetNode(key string) string {
	if len(r.keys) == 0 {
		return ""
	}
	hash := r.hash(key)
	idx := sort.Search(len(r.keys), func(i int) bool {
		return r.keys[i] >= hash
	})
	if idx == len(r.keys) {
		idx = 0 // Wrap around
	}
	return r.vnodeMap[r.keys[idx]]
}

func main() {
	// 虚拟节点设为100，以获得更好的均匀性
	ring := NewRing(100)

	// 场景一: 初始状态 (NodeA, NodeB, NodeC)
	fmt.Println("--- SCENE 1: Initial State (3 nodes) ---")
	nodes := []string{"NodeA", "NodeB", "NodeC"}
	for _, node := range nodes {
		ring.AddNode(node)
	}

	keys := []string{"Key 1", "Key 2", "Key 3", "Key 4", "Another Key", "Final Key"}
	distribution1 := make(map[string]string)
	fmt.Println("Initial data distribution:")
	for _, key := range keys {
		node := ring.GetNode(key)
		distribution1[key] = node
		fmt.Printf("  '%s' @%d° -> %s\n", key, ring.hash(key), node)
	}

	// 场景二: 扩容，增加 NodeD
	fmt.Println("\n--- SCENE 2: Add NodeD (Scale Up) ---")
	ring.AddNode("NodeD")

	distribution2 := make(map[string]string)
	rebalancedCount := 0
	fmt.Println("New data distribution:")
	for _, key := range keys {
		node := ring.GetNode(key)
		distribution2[key] = node
		if distribution1[key] != node {
			rebalancedCount++
			fmt.Printf("  [REBALANCE] '%s' moved from %s to -> %s\n", key, distribution1[key], node)
		} else {
			fmt.Printf("  '%s' -> %s (unchanged)\n", key, node)
		}
	}
	fmt.Printf(">>> After adding NodeD, %d out of %d keys were rebalanced.\n", rebalancedCount, len(keys))

	// 场景三: 缩容，移除 NodeB
	fmt.Println("\n--- SCENE 3: Remove NodeB (Scale Down) ---")
	ring.RemoveNode("NodeB")

	distribution3 := make(map[string]string)
	rebalancedCount = 0
	fmt.Println("Final data distribution:")
	for _, key := range keys {
		node := ring.GetNode(key)
		distribution3[key] = node
		if distribution2[key] != node {
			rebalancedCount++
			fmt.Printf("  [REBALANCE] '%s' moved from %s to -> %s\n", key, distribution2[key], node)
		} else {
			fmt.Printf("  '%s' -> %s (unchanged)\n", key, node)
		}
	}
	fmt.Printf(">>> After removing NodeB, %d out of %d keys were rebalanced.\n", rebalancedCount, len(keys))
}
