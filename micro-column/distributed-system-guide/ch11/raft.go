package main

import (
	"log"
	"sync"
)

// NodeState 定义节点的三种状态
type NodeState int

const (
	Follower NodeState = iota
	Candidate
	Leader
)

// LogEntry 定义日志条目的结构
type LogEntry struct {
	Term    int
	Command interface{}
}

// RaftNode 是 Raft 协议的核心实现
type RaftNode struct {
	mu          sync.Mutex
	id          int
	peers       []int
	state       NodeState
	currentTerm int
	votedFor    int

	// 日志
	log []LogEntry

	// volatile state on all servers
	commitIndex int
	lastApplied int

	// volatile state on leaders
	nextIndex  map[int]int
	matchIndex map[int]int

	// RPC 和 timer channels
	// (在真实实现中，这些会是网络连接)
	// 为简化，我们将在 main 函数中直接调用方法
}

// RequestVoteArgs 定义投票请求的结构
type RequestVoteArgs struct {
	Term         int
	CandidateID  int
	LastLogIndex int
	LastLogTerm  int
}

// RequestVoteResponse 定义投票响应
type RequestVoteResponse struct {
	Term        int
	VoteGranted bool
}

// AppendEntriesArgs 定义日志复制/心跳请求的结构
type AppendEntriesArgs struct {
	Term         int
	LeaderID     int
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []LogEntry
	LeaderCommit int
}

// AppendEntriesResponse 定义日志复制/心跳响应
type AppendEntriesResponse struct {
	Term    int
	Success bool
}

// RequestVote RPC handler
func (rf *RaftNode) RequestVote(args RequestVoteArgs, reply *RequestVoteResponse) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	// 规则1: 如果请求的任期小于当前任期，拒绝
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.VoteGranted = false
		return
	}

	// 如果收到更高任期的请求，转变为 Follower
	if args.Term > rf.currentTerm {
		rf.state = Follower
		rf.currentTerm = args.Term
		rf.votedFor = -1
	}

	reply.Term = rf.currentTerm

	// 规则2: 如果已投票或请求者的日志不够新，拒绝
	logIsUpToDate := (args.LastLogTerm > rf.lastLogTerm()) ||
		(args.LastLogTerm == rf.lastLogTerm() && args.LastLogIndex >= rf.lastLogIndex())

	if (rf.votedFor == -1 || rf.votedFor == args.CandidateID) && logIsUpToDate {
		rf.votedFor = args.CandidateID
		reply.VoteGranted = true
		log.Printf("Node %d: Voted for %d in Term %d.\n", rf.id, args.CandidateID, rf.currentTerm)
	} else {
		reply.VoteGranted = false
	}
}

// Helper functions to get last log's index and term
func (rf *RaftNode) lastLogIndex() int {
	return len(rf.log) - 1
}
func (rf *RaftNode) lastLogTerm() int {
	if len(rf.log) == 0 {
		return 0
	}
	return rf.log[len(rf.log)-1].Term
}

// AppendEntries RPC handler
func (rf *RaftNode) AppendEntries(args AppendEntriesArgs, reply *AppendEntriesResponse) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	// 规则1: 如果请求的任期小于当前任期，拒绝
	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.Success = false
		return
	}

	// 如果收到更高或相等任期的 Leader 的消息，确认自己的 Follower 身份
	if args.Term >= rf.currentTerm {
		rf.state = Follower
		rf.currentTerm = args.Term
		rf.votedFor = -1
	}

	reply.Term = rf.currentTerm

	// 规则2: 一致性检查
	// 如果 PrevLogIndex 处的日志不匹配，拒绝
	if rf.lastLogIndex() < args.PrevLogIndex {
		reply.Success = false
		return
	}
	if args.PrevLogIndex >= 0 && rf.log[args.PrevLogIndex].Term != args.PrevLogTerm {
		// 简化：在真实实现中，这里会截断日志
		reply.Success = false
		return
	}

	// 规则3 & 4: 追加新日志，并截断可能存在的冲突日志
	rf.log = append(rf.log[:args.PrevLogIndex+1], args.Entries...)
	reply.Success = true
	log.Printf("Node %d: Appended entries from Leader %d. Log length is now %d.\n", rf.id, args.LeaderID, len(rf.log))

	// 规则 5: 更新 commitIndex
	if args.LeaderCommit > rf.commitIndex {
		rf.commitIndex = min(args.LeaderCommit, rf.lastLogIndex())
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
