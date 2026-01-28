package main

import (
	"fmt"
	"sync"
)

// GCounter is a Grow-Only Counter CRDT
type GCounter struct {
	id      string
	payload map[string]int
	mu      sync.RWMutex
}

func NewGCounter(id string, nodeIDs []string) *GCounter {
	payload := make(map[string]int)
	for _, nodeID := range nodeIDs {
		payload[nodeID] = 0
	}
	return &GCounter{
		id:      id,
		payload: payload,
	}
}

// Increment the counter on the local node
func (c *GCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.payload[c.id]++
}

// Value returns the total value of the counter
func (c *GCounter) Value() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var total int
	for _, v := range c.payload {
		total += v
	}
	return total
}

// Merge combines this counter's state with another's
func (c *GCounter) Merge(other GCounter) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for id, value := range other.payload {
		if value > c.payload[id] {
			c.payload[id] = value
		}
	}
}

func main() {
	nodeIDs := []string{"A", "B", "C"}
	nodeA := NewGCounter("A", nodeIDs)
	nodeB := NewGCounter("B", nodeIDs)

	// 模拟 Node A 和 Node B 的并发、离线修改
	fmt.Println("--- Independent offline updates ---")
	nodeA.Increment()
	nodeA.Increment() // A 增加了 2
	nodeB.Increment() // B 增加了 1

	fmt.Printf("Node A: Value=%d, Payload=%v\n", nodeA.Value(), nodeA.payload)
	fmt.Printf("Node B: Value=%d, Payload=%v\n\n", nodeB.Value(), nodeB.payload)

	// 模拟 Node A 和 Node B 恢复在线并互相 Merge
	fmt.Println("--- Nodes come online and merge ---")

	// Merge A -> B
	fmt.Println("Merging A's state into B...")
	nodeB.Merge(nodeA)
	fmt.Printf("Node B after merge: Value=%d, Payload=%v\n", nodeB.Value(), nodeB.payload)

	// Merge B -> A
	fmt.Println("Merging B's state into A...")
	nodeA.Merge(nodeB)
	fmt.Printf("Node A after merge: Value=%d, Payload=%v\n\n", nodeA.Value(), nodeA.payload)

	fmt.Println("Final state is consistent on both nodes!")
}
