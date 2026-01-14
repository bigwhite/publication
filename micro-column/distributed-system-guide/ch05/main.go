package main

import "fmt"

// VectorClock 定义
type VectorClock map[string]int

// Compare 比较两个向量时钟
func (vc1 VectorClock) Compare(vc2 VectorClock) string {
	// 检查 vc1 是否 happen-before vc2
	isBefore := true
	for node, clock1 := range vc1 {
		if clock1 > vc2[node] {
			isBefore = false
			break
		}
	}
	if isBefore && len(vc1) < len(vc2) { // 检查是否有vc2中的key在vc1中不存在
		isBefore = true
	} else if isBefore && len(vc1) == len(vc2) {
		isSame := true
		for node, clock1 := range vc1 {
			if clock1 != vc2[node] {
				isSame = false
				break
			}
		}
		if isSame {
			isBefore = false
		}
	}

	// 检查 vc2 是否 happen-before vc1
	isAfter := true
	for node, clock2 := range vc2 {
		if clock2 > vc1[node] {
			isAfter = false
			break
		}
	}
	if isAfter && len(vc2) < len(vc1) {
		isAfter = true
	} else if isAfter && len(vc2) == len(vc1) {
		isSame := true
		for node, clock2 := range vc2 {
			if clock2 != vc1[node] {
				isSame = false
				break
			}
		}
		if isSame {
			isAfter = false
		}
	}

	if isBefore {
		return "HAPPENS_BEFORE"
	}
	if isAfter {
		return "HAPPENS_AFTER"
	}
	return "CONCURRENT"
}

// 模拟节点
type Node struct {
	ID    string
	Clock VectorClock
}

func NewNode(id string, peers []string) *Node {
	vc := make(VectorClock)
	for _, peer := range peers {
		vc[peer] = 0
	}
	return &Node{ID: id, Clock: vc}
}

func (n *Node) LocalEvent() {
	n.Clock[n.ID]++
	fmt.Printf("Event at %s, clock is now %v\n", n.ID, n.Clock)
}

func (n *Node) Send(receiver *Node) {
	n.Clock[n.ID]++
	fmt.Printf("%s sends message, clock is %v\n", n.ID, n.Clock)
	receiver.Receive(n.Clock)
}

func (n *Node) Receive(senderClock VectorClock) {
	// 更新本地时钟
	for id, clock := range senderClock {
		if clock > n.Clock[id] {
			n.Clock[id] = clock
		}
	}
	n.Clock[n.ID]++
	fmt.Printf("%s receives message, clock is now %v\n", n.ID, n.Clock)
}

func main() {
	peers := []string{"N1", "N2"}
	n1 := NewNode("N1", peers)
	n2 := NewNode("N2", peers)

	// 场景: N1 和 N2 并发地发生本地事件
	fmt.Println("--- SCENE 1: Concurrent Events ---")
	n1.LocalEvent() // eventA at N1
	eventA_clock := n1.Clock

	n2.LocalEvent() // eventB at N2
	eventB_clock := n2.Clock

	fmt.Printf("Comparing eventA %v and eventB %v: %s\n\n",
		eventA_clock, eventB_clock, eventA_clock.Compare(eventB_clock))

	// 场景: N1 发送消息给 N2, 建立因果关系
	fmt.Println("--- SCENE 2: Causal Events ---")
	n1.Send(n2)              // eventC (send at N1), eventD (receive at N2)
	eventC_clock := n1.Clock // Note: this is after send, so it's a new state
	eventD_clock := n2.Clock

	fmt.Printf("Comparing send event's cause %v and receive event %v: %s\n\n",
		eventC_clock, eventD_clock, eventC_clock.Compare(eventD_clock))

}
