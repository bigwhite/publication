package main

import (
	"fmt"
	"log"
	"sync"
)

// TCCSubTransaction 定义TCC子事务接口
type TCCSubTransaction interface {
	Try() error
	Confirm()
	Cancel()
}

// ---- 各个服务的TCC实现 ----
var accountBalance = 1000
var frozenBalance = 0
var accountMu sync.Mutex

type AccountTX struct{ amount int }

func (tx *AccountTX) Try() error {
	accountMu.Lock()
	defer accountMu.Unlock()
	log.Printf("TRY: Account Service - Current balance: %d. Trying to freeze %d...\n", accountBalance, tx.amount)
	if accountBalance < tx.amount {
		return fmt.Errorf("insufficient account balance")
	}
	accountBalance -= tx.amount
	frozenBalance += tx.amount
	log.Printf("...Success. Frozen balance: %d, Available balance: %d\n", frozenBalance, accountBalance)
	return nil
}
func (tx *AccountTX) Confirm() {
	accountMu.Lock()
	defer accountMu.Unlock()
	log.Println("CONFIRM: Account Service - Deducting frozen balance...")
	frozenBalance -= tx.amount
	log.Printf("...Success. Final balance: %d\n", accountBalance)
}
func (tx *AccountTX) Cancel() {
	accountMu.Lock()
	defer accountMu.Unlock()
	log.Println("CANCEL: Account Service - Unfreezing balance...")
	frozenBalance -= tx.amount
	accountBalance += tx.amount
	log.Printf("...Success. Rolled back balance to %d\n", accountBalance)
}

type RedPacketTX struct{}

func (tx *RedPacketTX) Try() error {
	log.Println("TRY: RedPacket Service - Checking packet status...")
	// 模拟红包已过期
	return fmt.Errorf("red packet expired")
}
func (tx *RedPacketTX) Confirm() { log.Println("CONFIRM: RedPacket Service - Using red packet...") }
func (tx *RedPacketTX) Cancel()  { log.Println("CANCEL: RedPacket Service - Releasing red packet...") }

// TCCCoordinator: TCC协调器
type TCCCoordinator struct{ transactions []TCCSubTransaction }

func (c *TCCCoordinator) Add(tx TCCSubTransaction) { c.transactions = append(c.transactions, tx) }
func (c *TCCCoordinator) Execute() {
	log.Println("--- Starting TCC Transaction (TRY Phase) ---")

	var successfulTries []TCCSubTransaction
	tryFailed := false
	for i, tx := range c.transactions {
		if err := tx.Try(); err != nil {
			log.Printf("!!! ERROR on TRY step %d: %v. Starting Cancel Phase...\n", i+1, err)
			tryFailed = true
			break
		}
		successfulTries = append(successfulTries, tx)
	}

	if tryFailed {
		log.Println("\n--- TCC Transaction Failed (CANCEL Phase) ---")
		for _, tx := range successfulTries {
			tx.Cancel()
		}
	} else {
		log.Println("\n--- TCC Transaction Succeeded (CONFIRM Phase) ---")
		for _, tx := range c.transactions {
			tx.Confirm()
		}
	}
	log.Println("--- TCC Transaction Finished ---")
}

func main() {
	tcc := &TCCCoordinator{}
	tcc.Add(&AccountTX{amount: 100})
	tcc.Add(&RedPacketTX{})
	tcc.Execute()
}
