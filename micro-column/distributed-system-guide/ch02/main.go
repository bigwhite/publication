package main

import (
	"fmt"
	"sync"
)

// 一个简化的银行服务
type BankService struct {
	mu            sync.Mutex
	balances      map[string]int
	processedReqs map[string]bool // 存储已处理的请求ID
}

func NewBankService() *BankService {
	return &BankService{
		balances:      map[string]int{"Alice": 1000, "Bob": 1000},
		processedReqs: make(map[string]bool),
	}
}

// 转账操作，增加了幂等性检查
func (s *BankService) Transfer(reqID, from, to string, amount int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 哲学：假设这个请求可能被重试，先检查我们是否已经处理过它。
	if s.processedReqs[reqID] {
		fmt.Printf("Request %s already processed, skipping.\n", reqID)
		return nil // 幂等处理：直接返回成功
	}

	if s.balances[from] < amount {
		return fmt.Errorf("insufficient funds")
	}

	s.balances[from] -= amount
	s.balances[to] += amount
	s.processedReqs[reqID] = true // 标记为已处理

	fmt.Printf("Processed request %s: %s -> %s, amount %d. New balances: %v\n", reqID, from, to, amount, s.balances)
	return nil
}

func main() {
	service := NewBankService()
	reqID := "tx-12345" // 每个业务操作都应该有一个唯一的ID

	// 模拟一次网络调用
	service.Transfer(reqID, "Alice", "Bob", 100)

	// 模拟因为网络超时，客户端发起了重试
	fmt.Println("\n--- Client retries due to timeout ---")
	service.Transfer(reqID, "Alice", "Bob", 100)
}
