package main

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
)

// HashRing 一致性哈希环
type HashRing struct {
	// sortedKeys 存储所有虚拟节点的哈希值，必须保持有序，以便二分查找
	sortedKeys []int
	// virtualMap 存储 虚拟节点哈希 -> 真实节点名称
	virtualMap map[int]string
	// replicas 每个真实节点对应的虚拟节点数量（为了数据平衡）
	replicas int
}

func NewHashRing(replicas int) *HashRing {
	return &HashRing{
		replicas:   replicas,
		virtualMap: make(map[int]string),
	}
}

// AddNode 添加一个真实节点
func (h *HashRing) AddNode(nodeName string) {
	for i := 0; i < h.replicas; i++ {
		// 生成虚拟节点名称，例如 "NodeA#0", "NodeA#1"
		virtualKey := nodeName + "#" + strconv.Itoa(i)
		// 计算哈希值
		hash := int(crc32.ChecksumIEEE([]byte(virtualKey)))

		h.sortedKeys = append(h.sortedKeys, hash)
		h.virtualMap[hash] = nodeName
	}
	// 每次添加后，必须重新排序，保证二分查找的正确性
	sort.Ints(h.sortedKeys)
}

// GetNode 根据数据 Key 获取对应的真实节点
func (h *HashRing) GetNode(key string) string {
	if len(h.sortedKeys) == 0 {
		return ""
	}

	hash := int(crc32.ChecksumIEEE([]byte(key)))

	// 核心：二分查找
	// sort.Search(n, f) 返回 [0, n) 中满足 f(i) 为 true 的最小下标 i
	// 我们要找第一个 nodeHash >= keyHash
	idx := sort.Search(len(h.sortedKeys), func(i int) bool {
		return h.sortedKeys[i] >= hash
	})

	// 如果 idx == len，说明 keyHash 比环上所有节点都大
	// 根据环形逻辑，它应该归属于第一个节点（绕回起点）
	if idx == len(h.sortedKeys) {
		idx = 0
	}

	// 映射回真实节点
	return h.virtualMap[h.sortedKeys[idx]]
}

func main() {
	ring := NewHashRing(3) // 每个节点生成 3 个虚拟节点

	// 添加 3 台服务器
	ring.AddNode("Server_A")
	ring.AddNode("Server_B")
	ring.AddNode("Server_C")

	fmt.Println("--- Hash Ring Nodes ---")
	for _, k := range ring.sortedKeys {
		fmt.Printf("Hash: %d -> Node: %s\n", k, ring.virtualMap[k])
	}

	fmt.Println("\n--- Key Routing ---")
	keys := []string{"User_1", "User_2", "Order_123", "Session_ABC"}
	for _, k := range keys {
		node := ring.GetNode(k)
		fmt.Printf("Key: %-12s => Node: %s\n", k, node)
	}
}
