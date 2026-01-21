package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// 创建一个3节点的集群
	cluster := make([]*RaftNode, 3)
	peers := []int{0, 1, 2}
	for i := 0; i < 3; i++ {
		cluster[i] = &RaftNode{
			id:          i,
			peers:       peers,
			state:       Follower,
			currentTerm: 0,
			votedFor:    -1,
			log:         []LogEntry{},
		}
	}

	log.Println("--- Raft Cluster Simulation Start (3 Nodes) ---")

	// --- 模拟选举 ---
	log.Println("\n--- SCENE 1: ELECTION ---")
	// 假设 Node 0 选举超时，成为 Candidate
	candidate := cluster[0]
	candidate.mu.Lock()
	candidate.state = Candidate
	candidate.currentTerm = 1
	candidate.votedFor = 0
	candidate.mu.Unlock()
	log.Printf("Node 0 becomes Candidate for Term 1.\n")

	votes := 1
	var wg sync.WaitGroup
	for _, peerID := range candidate.peers {
		if peerID == candidate.id {
			continue
		}

		wg.Add(1)
		go func(peer *RaftNode) {
			defer wg.Done()
			args := RequestVoteArgs{
				Term:         candidate.currentTerm,
				CandidateID:  candidate.id,
				LastLogIndex: candidate.lastLogIndex(),
				LastLogTerm:  candidate.lastLogTerm(),
			}
			var reply RequestVoteResponse
			peer.RequestVote(args, &reply)

			if reply.VoteGranted {
				votes++
			}
		}(cluster[peerID])
	}
	wg.Wait()

	if votes > len(cluster)/2 {
		log.Printf("Candidate 0 received %d votes and becomes LEADER for Term 1!\n", votes)
		candidate.state = Leader
	} else {
		log.Printf("Candidate 0 failed to become leader.\n")
	}

	// --- 模拟日志复制 ---
	log.Println("\n--- SCENE 2: LOG REPLICATION ---")
	leader := candidate

	// Leader 收到一条新指令
	newEntry := LogEntry{Term: leader.currentTerm, Command: "SET x = 10"}
	leader.mu.Lock()
	leader.log = append(leader.log, newEntry)
	leader.mu.Unlock()
	log.Printf("Leader 0: Appended new entry {T:%d, C:%s} to its log.\n", newEntry.Term, newEntry.Command)

	// Leader 将新日志复制给 Followers
	replicatedCount := 1
	for _, peerID := range leader.peers {
		if peerID == leader.id {
			continue
		}

		wg.Add(1)
		go func(peer *RaftNode) {
			defer wg.Done()
			args := AppendEntriesArgs{
				Term:         leader.currentTerm,
				LeaderID:     leader.id,
				PrevLogIndex: leader.lastLogIndex() - 1,
				PrevLogTerm:  0, // 假设之前日志为空
				Entries:      []LogEntry{newEntry},
			}
			var reply AppendEntriesResponse
			peer.AppendEntries(args, &reply)

			if reply.Success {
				replicatedCount++
			}
		}(cluster[peerID])
	}
	wg.Wait()

	if replicatedCount > len(cluster)/2 {
		leader.commitIndex = leader.lastLogIndex()
		log.Printf("SUCCESS: Entry has been replicated to a majority (%d/%d). Leader commits index %d.\n", replicatedCount, len(cluster), leader.commitIndex)
	} else {
		log.Println("FAILURE: Entry failed to replicate to a majority.")
	}
}
