package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)

	// 在这10s期间，我们多次触发SIGINT信号
	time.Sleep(10 * time.Second)

	for {
		select {
		case s := <-c:
			fmt.Println("c: 获取异步信号", s)
		default:
			fmt.Println("c: 没有信号, 退出")
			return
		}
	}
}
