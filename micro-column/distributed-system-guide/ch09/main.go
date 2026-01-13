// ch09/main.go
package main

import (
	"log"
	"sync"
	"time"
)

// --- Acceptor 侧 ---
type Acceptor struct {
	id          int
	promisedID  int
	acceptedID  int
	acceptedVal string
	mu          sync.Mutex
}

func (a *Acceptor) Prepare(proposalID int, promiseChan chan<- Promise) {
	a.mu.Lock()
	defer a.mu.Unlock()

	log.Printf("Acceptor %d: Received Prepare(ID=%d). My promisedID is %d.", a.id, proposalID, a.promisedID)
	if proposalID > a.promisedID {
		a.promisedID = proposalID
		log.Printf(" -> Promised ID %d.", proposalID)
		promiseChan <- Promise{
			Promised:     true,
			AcceptedID:   a.acceptedID,
			AcceptedVal:  a.acceptedVal,
			FromAcceptor: a.id,
		}
	} else {
		log.Printf(" -> Rejected.")
		promiseChan <- Promise{Promised: false, FromAcceptor: a.id}
	}
}

func (a *Acceptor) Accept(proposalID int, value string, acceptChan chan<- bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	log.Printf("Acceptor %d: Received Accept(ID=%d, Val='%s'). My promisedID is %d.", a.id, proposalID, value, a.promisedID)
	if proposalID >= a.promisedID {
		a.promisedID = proposalID
		a.acceptedID = proposalID
		a.acceptedVal = value
		log.Printf(" -> Accepted.")
		acceptChan <- true
	} else {
		log.Printf(" -> Rejected.")
		acceptChan <- false
	}
}

// --- Proposer 侧 ---
type Promise struct {
	Promised     bool
	AcceptedID   int
	AcceptedVal  string
	FromAcceptor int
}

type Proposer struct {
	id         string
	proposalID int
	acceptors  []*Acceptor
}

func (p *Proposer) Propose(value string) (string, bool) {
	log.Printf("\n--- Proposer %s starts proposing '%s' with ID %d ---\n", p.id, value, p.proposalID)
	// Phase 1: Prepare
	promiseChan := make(chan Promise, len(p.acceptors))
	for _, a := range p.acceptors {
		go a.Prepare(p.proposalID, promiseChan)
	}

	quorumSize := len(p.acceptors)/2 + 1
	promises := make([]Promise, 0)
	for i := 0; i < len(p.acceptors); i++ {
		promise := <-promiseChan
		if promise.Promised {
			promises = append(promises, promise)
		}
	}

	if len(promises) < quorumSize {
		log.Printf("Proposer %s: Prepare phase failed. Not enough promises (%d/%d).\n", p.id, len(promises), quorumSize)
		return "", false
	}
	log.Printf("Proposer %s: Prepare phase successful. Received %d promises.\n", p.id, len(promises))

	// 决定要 Accept 的值
	highestAcceptedID := -1
	valueToPropose := value // 默认用自己的值
	hasPrevValue := false
	for _, promise := range promises {
		// 只有当 Acceptor 真正接受过一个值 (AcceptedID > 0)，才认为它的历史有效
		if promise.AcceptedID > 0 && promise.AcceptedID > highestAcceptedID {
			highestAcceptedID = promise.AcceptedID
			valueToPropose = promise.AcceptedVal
			hasPrevValue = true
		}
	}
	if hasPrevValue {
		log.Printf("Proposer %s: Found a previously accepted value ('%s') with ID %d. Will propose this value instead.\n", p.id, valueToPropose, highestAcceptedID)
	} else {
		log.Printf("Proposer %s: No previously accepted value found. Proposing my own value ('%s').\n", p.id, valueToPropose)
	}

	// Phase 2: Accept
	acceptChan := make(chan bool, len(p.acceptors))
	// 只向做出承诺的 Acceptor 发送 Accept 请求
	for _, a := range p.acceptors {
		go a.Accept(p.proposalID, valueToPropose, acceptChan)
	}

	acceptCount := 0
	for i := 0; i < len(p.acceptors); i++ {
		if <-acceptChan {
			acceptCount++
		}
	}

	if acceptCount >= quorumSize {
		log.Printf("Proposer %s: SUCCESS! Value '%s' has been CHOSEN by consensus.\n", p.id, valueToPropose)
		return valueToPropose, true
	}

	log.Printf("Proposer %s: Accept phase failed. Not enough accepts (%d/%d).\n", p.id, acceptCount, quorumSize)
	return "", false
}

func main() {
	// 创建3个 Acceptor
	acceptors := []*Acceptor{
		{id: 1},
		{id: 2},
		{id: 3},
	}

	// 场景一：Proposer P1 首次提出 'Value A'
	p1 := &Proposer{id: "P1", proposalID: 10, acceptors: acceptors}
	p1.Propose("Value A")

	// 等待 P1 的 goroutines 结束, 确保状态被更新
	time.Sleep(100 * time.Millisecond)

	// 场景二：Proposer P2 提出 'Value B'，它的提案ID更高
	p2 := &Proposer{id: "P2", proposalID: 20, acceptors: acceptors}
	p2.Propose("Value B")
}
