package main

import (
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	println("main goroutine: start to work...")
	go func() {
		println("goroutine1: start to work...")
		time.Sleep(5 * time.Second)
		println("goroutine1: work done!")
		wg.Done()
	}()

	go func() {
		println("goroutine2: start to work...")
		time.Sleep(1 * time.Second)
		panic("division by zero")
		println("goroutine2: work done!")
		wg.Done()
	}()

	wg.Wait()
	println("main goroutine: work done!")
}
