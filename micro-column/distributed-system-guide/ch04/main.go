package main

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// VersionedData 包含数据和版本信息
type VersionedData struct {
	Value   string
	Version int64 // 使用纳秒时间戳作为版本
}

// Node 模拟一个副本节点
type Node struct {
	id   string
	data map[string]VersionedData
	mu   sync.RWMutex
}

func (n *Node) Write(key, value string, version int64) {
	n.mu.Lock()
	defer n.mu.Unlock()
	// LWW: 只有当新版本更高时才写入
	if n.data[key].Version < version {
		n.data[key] = VersionedData{Value: value, Version: version}
		fmt.Printf("Node %s: Wrote key '%s' with value '%s' (version %d)\n", n.id, key, value, version)
	}
}

func (n *Node) Read(key string) VersionedData {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.data[key]
}

// Cluster 模拟无主集群
type Cluster struct {
	nodes []*Node
	N     int
	W     int
	R     int
}

// Write 操作，实现了W Quorum
func (c *Cluster) Write(key, value string) {
	version := time.Now().UnixNano() // 使用时间戳作为版本号 (LWW)
	acks := make(chan bool, c.N)

	for _, node := range c.nodes {
		go func(n *Node) {
			n.Write(key, value, version)
			acks <- true
		}(node)
	}

	// 等待W个确认
	for i := 0; i < c.W; i++ {
		<-acks
	}
	fmt.Printf(">>> Write successful for key '%s' with W=%d acks.\n\n", key, c.W)
}

// ReadWithRepair 操作，实现了R Quorum和读修复
func (c *Cluster) ReadWithRepair(key string) string {
	results := make(chan VersionedData, c.N)
	// 为了演示读修复，我们从所有节点读取
	readCount := len(c.nodes)
	for _, node := range c.nodes {
		go func(n *Node) {
			results <- n.Read(key)
		}(node)
	}

	var receivedResults []VersionedData
	for i := 0; i < readCount; i++ {
		receivedResults = append(receivedResults, <-results)
	}

	// 即使我们读取了所有节点，也只关心R个最快的响应（这里为了简化，我们处理所有）
	// 在真实世界中，客户端会等待R个响应即可

	// 找出最新的版本
	sort.Slice(receivedResults, func(i, j int) bool {
		return receivedResults[i].Version > receivedResults[j].Version
	})
	latest := receivedResults[0]
	fmt.Printf(">>> Read for key '%s' got %d results. Latest version is %d ('%s').\n", key, len(receivedResults), latest.Version, latest.Value)

	// **读修复 (Read Repair)**
	go func() {
		for _, node := range c.nodes {
			if node.Read(key).Version < latest.Version {
				fmt.Printf("[Read Repair] Found stale node %s. Updating it to version %d.\n", node.id, latest.Version)
				node.Write(key, latest.Value, latest.Version)
			}
		}
		fmt.Println("[Read Repair] Repair process completed.")
	}()

	return latest.Value
}

func main() {
	// 初始化集群 N=3, W=2, R=2
	nodes := []*Node{
		{id: "N1", data: make(map[string]VersionedData)},
		{id: "N2", data: make(map[string]VersionedData)},
		{id: "N3", data: make(map[string]VersionedData)},
	}
	cluster := &Cluster{nodes: nodes, N: 3, W: 2, R: 2}

	// 场景1: 模拟一个节点N3宕机，写入依然成功
	fmt.Println("--- SCENE 1: Write with one node down (simulated) ---")
	// 临时移除 N3
	activeNodes := []*Node{nodes[0], nodes[1]}
	tempCluster := &Cluster{nodes: activeNodes, N: 2, W: 2, R: 2}
	tempCluster.Write("greeting", "hello")

	// 恢复N3, 此时N3没有'greeting'的数据
	fmt.Printf("Node N3 is back online, but it missed the write for 'greeting'.\n\n")

	// 场景2: 读取'greeting'，触发读修复
	fmt.Println("--- SCENE 2: Read triggers Read Repair ---")
	cluster.ReadWithRepair("greeting")

	// 等待读修复goroutine完成
	time.Sleep(100 * time.Millisecond)

	// 验证N3是否已被修复
	fmt.Println("\n--- Verifying repair on N3 ---")
	dataOnN3 := nodes[2].Read("greeting")
	fmt.Printf("Data on N3 for key 'greeting' is now: '%s' (version %d)\n", dataOnN3.Value, dataOnN3.Version)
}
