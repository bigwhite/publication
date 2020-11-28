package main

import "sync"

func foo(mu sync.Mutex) {
	mu.Lock()
	mu.Unlock()
}

func main() {
	var mu sync.Mutex
	foo(mu)
	mu1 := mu
	mu1.Lock()
	mu1.Unlock()

	pMu := &mu
	pMu1 := &mu1
	*pMu1 = *pMu
}
