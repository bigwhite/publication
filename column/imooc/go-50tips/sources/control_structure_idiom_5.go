package main

import (
	"fmt"
	"time"
)

func recvFromUnbufferedChannel() {
	var c = make(chan int)

	go func() {
		time.Sleep(time.Second * 3)
		c <- 1
		c <- 2
		c <- 3
		close(c)
	}()

	for v := range c {
		fmt.Println(v)
	}
}

func recvFromNilChannel() {
	var c chan int

	// 程序将一直阻塞在这里
	for v := range c {
		fmt.Println(v)
	}

}

func main() {
	recvFromUnbufferedChannel()
	recvFromNilChannel()
}
