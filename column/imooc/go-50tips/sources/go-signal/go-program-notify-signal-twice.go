package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	c1 := make(chan os.Signal, 1)
	c2 := make(chan os.Signal, 1)

	signal.Notify(c1, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(c2, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s := <-c1
		fmt.Println("c1: 收到异步信号", s)
	}()

	s := <-c2
	fmt.Println("c2: 收到异步信号", s)
	time.Sleep(5 * time.Second)
}
