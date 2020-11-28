package main

import (
	"fmt"
	"sync"
	"time"
)

func safeRun(g func()) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("caught a panic:", e)
		}
	}()

	g()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	println("main goroutine: start to work...")
	go safeRun(func() {
		defer wg.Done()
		println("goroutine1: start to work...")
		time.Sleep(5 * time.Second)
		println("goroutine1: work done!")
	})

	go safeRun(func() {
		defer wg.Done()
		println("goroutine2: start to work...")
		time.Sleep(1 * time.Second)
		panic("division by zero")
		println("goroutine2: work done!")
	})

	wg.Wait()
	println("main goroutine: work done!")
}
