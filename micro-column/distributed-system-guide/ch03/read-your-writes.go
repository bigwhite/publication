package main

import (
	"fmt"
	"sync"
	"time"
)

// 一个简化的数据存储
type Store struct {
	mu    sync.RWMutex
	data  map[string]string
	lag   time.Duration // 模拟复制延迟
	epoch int           // 数据的版本
}

func (s *Store) Write(key, value string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	time.Sleep(10 * time.Millisecond) // 模拟写入耗时
	s.data[key] = value
	s.epoch++
	return s.epoch
}

func (s *Store) Read(key string) (string, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// 模拟延迟
	time.Sleep(s.lag)
	return s.data[key], s.epoch
}

// 模拟我们的主从集群
var leader = &Store{data: make(map[string]string)}
var follower = &Store{data: make(map[string]string), lag: 100 * time.Millisecond} // 100ms 延迟

// 异步复制过程
func replicate(key, value string) {
	go func() {
		time.Sleep(follower.lag) // 模拟网络传输延迟
		follower.Write(key, value)
		fmt.Printf("[Replication] Follower updated to epoch %d\n", follower.epoch)
	}()
}

// 模拟用户最近有写入的缓存
var recentWriters = make(map[string]time.Time)
var rwMutex sync.Mutex

func SmartRead(user, key string) (string, int) {
	rwMutex.Lock()
	lastWriteTime, hasWritten := recentWriters[user]
	rwMutex.Unlock()

	// 哲学：如果用户最近写入过，为了保证体验，从Leader读
	if hasWritten && time.Since(lastWriteTime) < 1*time.Minute {
		fmt.Printf("[Read Router] User '%s' wrote recently. Reading from LEADER.\n", user)
		return leader.Read(key)
	}

	fmt.Printf("[Read Router] User '%s' has not written recently. Reading from FOLLOWER.\n", user)
	return follower.Read(key)
}

func main() {
	user := "gopher"
	key := "profile_status"
	value := "Coding distributed systems!"

	// 1. 用户写入
	fmt.Println("--- User writes data ---")
	epoch := leader.Write(key, value)
	replicate(key, value) // 异步复制

	rwMutex.Lock()
	recentWriters[user] = time.Now()
	rwMutex.Unlock()
	fmt.Printf("Leader accepted write. Data is '%s', epoch is %d\n\n", value, epoch)

	time.Sleep(50 * time.Millisecond) // 等待一小会儿，但小于复制延迟

	// 2. 用户立即读取
	fmt.Println("--- User reads immediately ---")
	readValue, readEpoch := SmartRead(user, key)
	fmt.Printf("User '%s' read value: '%s' (from epoch %d)\n\n", user, readValue, readEpoch)

	// 3. 另一个用户读取
	fmt.Println("--- Another user reads ---")
	otherUser := "rustacean"
	readValue, readEpoch = SmartRead(otherUser, key)
	fmt.Printf("User '%s' read value: '%s' (from epoch %d)\n", otherUser, readValue, readEpoch)

	// 等待复制完成
	time.Sleep(100 * time.Millisecond)
	fmt.Println("\n--- After replication completes ---")
	readValue, readEpoch = SmartRead(otherUser, key)
	fmt.Printf("User '%s' reads again, value: '%s' (from epoch %d)\n", otherUser, readValue, readEpoch)
}
