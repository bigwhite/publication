package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Participant 模拟一个参与者
type Participant struct {
	id      string
	vote    chan bool // 用于投票
	ack     chan bool // 用于确认
	isReady bool      // 是否已准备好
}

func (p *Participant) Prepare() {
	// 模拟检查和写日志
	log.Printf("Participant %s: Preparing... (writing to WAL, locking resources)\n", p.id)
	time.Sleep(50 * time.Millisecond)
	p.isReady = true
	log.Printf("Participant %s: VOTE-COMMIT\n", p.id)
	p.vote <- true // 投票同意
}

func (p *Participant) Commit() {
	if p.isReady {
		log.Printf("Participant %s: Committing transaction, releasing locks.\n", p.id)
		p.ack <- true
	}
}

func (p *Participant) Rollback() {
	if p.isReady {
		log.Printf("Participant %s: Rolling back transaction, releasing locks.\n", p.id)
		p.isReady = false
		p.ack <- true
	}
}

// Coordinator 模拟协调者
type Coordinator struct {
	participants []*Participant
}

func (c *Coordinator) RunTransaction() {
	log.Println("--- PHASE 1: PREPARE ---")

	// 1. 发送 PREPARE 请求
	for _, p := range c.participants {
		go p.Prepare()
	}

	// 2. 收集投票
	allVoteCommit := true
	for _, p := range c.participants {
		select {
		case vote := <-p.vote:
			if !vote {
				allVoteCommit = false
				log.Printf("Coordinator: Received VOTE-ABORT from a participant.\n")
			}
		case <-time.After(100 * time.Millisecond):
			allVoteCommit = false
			log.Printf("Coordinator: Timeout waiting for a vote.\n")
		}
	}

	// *** 模拟协调者在此刻崩溃！ ***
	if simulateCrash {
		log.Fatalln("!!!!!! COORDINATOR CRASHES AFTER PHASE 1, BEFORE PHASE 2 !!!!!!")
		return // 进程退出
	}

	log.Println("\n--- PHASE 2: COMMIT/ROLLBACK ---")

	// 3. 做出决定并广播
	if allVoteCommit {
		log.Println("Coordinator: All participants voted COMMIT. Broadcasting COMMIT...")
		for _, p := range c.participants {
			go p.Commit()
		}
	} else {
		log.Println("Coordinator: Found ABORT vote. Broadcasting ROLLBACK...")
		for _, p := range c.participants {
			go p.Rollback()
		}
	}

	// 4. 等待ACK
	for _, p := range c.participants {
		<-p.ack
	}

	log.Println("\n--- TRANSACTION COMPLETE ---")
}

var simulateCrash = false

func main() {
	p1 := &Participant{id: "P1", vote: make(chan bool, 1), ack: make(chan bool, 1)}
	p2 := &Participant{id: "P2", vote: make(chan bool, 1), ack: make(chan bool, 1)}
	coordinator := &Coordinator{participants: []*Participant{p1, p2}}

	fmt.Println("--- Running normal transaction ---")
	coordinator.RunTransaction()

	fmt.Println("\n\n--- Running transaction with coordinator crash ---")
	simulateCrash = true
	// 在一个goroutine中运行，因为log.Fatalln会终止程序
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// 模拟参与者状态
				log.Println("Coordinator is dead. What should Participants do?")
				log.Printf("Participant %s is now FROZEN (isReady=%v), waiting for coordinator...\n", p1.id, p1.isReady)
				log.Printf("Participant %2s is now FROZEN (isReady=%v), waiting for coordinator...\n", p2.id, p2.isReady)
				wg.Done()
			}
		}()
		coordinator.RunTransaction()
	}()
	wg.Wait()
}
