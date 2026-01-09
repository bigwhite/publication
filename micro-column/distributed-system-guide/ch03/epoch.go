package main

import (
	"fmt"
	"sync"
)

type Node struct {
	id           string
	currentEpoch int
	mu           sync.Mutex
}

// HandleWriteRequest 模拟节点处理来自Leader的写入请求
func (n *Node) HandleWriteRequest(leaderID string, leaderEpoch int, data string) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	fmt.Printf("Node %s: Received write request from %s (Epoch %d). My current epoch is %d.\n", n.id, leaderID, leaderEpoch, n.currentEpoch)

	// 哲学：只听从当前或更新“朝代”的命令
	if leaderEpoch < n.currentEpoch {
		return fmt.Errorf("rejected: leader's epoch %d is old, my epoch is %d", leaderEpoch, n.currentEpoch)
	}

	// 如果leader的epoch更新，则更新自己的epoch并接受写入
	if leaderEpoch > n.currentEpoch {
		fmt.Printf("Node %s: Saw a new epoch %d, updating from %d.\n", n.id, leaderEpoch, n.currentEpoch)
		n.currentEpoch = leaderEpoch
	}

	fmt.Printf("Node %s: Accepted write '%s' from %s (Epoch %d).\n", n.id, data, leaderID, n.currentEpoch)
	// ... 实际的写入逻辑 ...
	return nil
}

// SimulateNewElection 模拟一次新的选举，递增epoch
func SimulateNewElection(nodes []*Node, newLeader *Node) {
	highestEpoch := 0
	for _, n := range nodes {
		if n.currentEpoch > highestEpoch {
			highestEpoch = n.currentEpoch
		}
	}
	newEpoch := highestEpoch + 1
	fmt.Printf("\n--- NEW ELECTION! New Epoch will be %d. %s is the new leader. ---\n", newEpoch, newLeader.id)

	for _, n := range nodes {
		n.currentEpoch = newEpoch
	}
}

func main() {
	// 初始化集群，初始epoch为1
	nodes := []*Node{
		{id: "L1 (Old Leader)", currentEpoch: 1},
		{id: "F2", currentEpoch: 1},
		{id: "F3", currentEpoch: 1},
	}

	oldLeader := nodes[0]
	follower := nodes[1]

	// 场景1: 正常写入
	fmt.Println("--- Normal Operation ---")
	follower.HandleWriteRequest(oldLeader.id, oldLeader.currentEpoch, "data A")

	// 场景2: 发生网络分区，F2和F3选举出了新Leader（我们手动模拟这个过程）
	// 注意此时oldLeader L1还不知道自己被罢免了，它的epoch还是1
	newLeader := follower // F2成为新Leader
	SimulateNewElection(nodes, newLeader)

	// 场景3: 新Leader L2 (原F2) 正常写入
	fmt.Println("\n--- New Leader Operation ---")
	follower.HandleWriteRequest(newLeader.id, newLeader.currentEpoch, "data B")
	// 另一个节点 F3 也能正常处理
	nodes[2].HandleWriteRequest(newLeader.id, newLeader.currentEpoch, "data C")

	// 场景4: 网络分区恢复，被罢免的旧Leader L1苏醒，试图发出旧epoch的命令
	fmt.Println("\n--- Deposed Leader Tries to Write ---")
	err := follower.HandleWriteRequest(oldLeader.id, oldLeader.currentEpoch, "data Z (stale write)")
	if err != nil {
		fmt.Printf("SUCCESS: Follower correctly rejected the stale write from %s. Reason: %v\n", oldLeader.id, err)
	}
}
