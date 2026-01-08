package main

import (
	"fmt"
	"log"
)

// SagaSubTransaction 定义子事务接口
type SagaSubTransaction interface {
	Execute() error
	Compensate() error
}

// ---- 各个服务的子事务实现 ----
type OrderTX struct{ orderID string }

func (tx *OrderTX) Execute() error {
	log.Printf("Step 1: Order service - Creating order '%s'...\n", tx.orderID)
	return nil
}
func (tx *OrderTX) Compensate() error {
	log.Printf("COMPENSATE: Order service - Canceling order '%s'...\n", tx.orderID)
	return nil
}

type InventoryTX struct{ productID string }

func (tx *InventoryTX) Execute() error {
	log.Printf("Step 2: Inventory service - Deducting stock for product '%s'...\n", tx.productID)
	return nil
}
func (tx *InventoryTX) Compensate() error {
	log.Printf("COMPENSATE: Inventory service - Adding stock back for product '%s'...\n", tx.productID)
	return nil
}

type PaymentTX struct {
	userID string
	amount float64
}

func (tx *PaymentTX) Execute() error {
	log.Printf("Step 3: Payment service - Processing payment for user '%s'...\n", tx.userID)
	return fmt.Errorf("insufficient funds")
}
func (tx *PaymentTX) Compensate() error {
	log.Printf("COMPENSATE: Payment service - Refunding user '%s'...\n", tx.userID)
	return nil
}

// SagaOrchestrator: SAGA 编排器
type SagaOrchestrator struct{ transactions []SagaSubTransaction }

func (s *SagaOrchestrator) Add(tx SagaSubTransaction) { s.transactions = append(s.transactions, tx) }
func (s *SagaOrchestrator) Execute() {
	log.Println("--- Starting SAGA Transaction ---")
	for i, tx := range s.transactions {
		if err := tx.Execute(); err != nil {
			log.Printf("!!! ERROR on step %d: %v. Starting compensation...\n", i+1, err)
			s.compensate(i - 1)
			return
		}
	}
	log.Println("--- SAGA Transaction Completed Successfully ---")
}
func (s *SagaOrchestrator) compensate(fromIndex int) {
	for i := fromIndex; i >= 0; i-- {
		if err := s.transactions[i].Compensate(); err != nil {
			log.Fatalf("FATAL: Compensation for step %d failed: %v.", i+1, err)
		}
	}
	log.Println("--- SAGA Compensation Completed ---")
}

func main() {
	saga := &SagaOrchestrator{}
	saga.Add(&OrderTX{orderID: "order-123"})
	saga.Add(&InventoryTX{productID: "prod-456"})
	saga.Add(&PaymentTX{userID: "user-789", amount: 99.99})
	saga.Execute()
}
