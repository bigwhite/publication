package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	var wg sync.WaitGroup
	c := make(chan os.Signal, 2)

	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGINT)

	// Block until any signal is received.
	wg.Add(1)
	go func() {
		for {
			s := <-c
			fmt.Println("c: 收到异步信号", s)
		}
		wg.Done()
	}()
	wg.Wait()
}
