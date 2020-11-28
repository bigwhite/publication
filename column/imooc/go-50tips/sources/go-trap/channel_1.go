package main

import "fmt"

func main() {
	var nilChan chan int
	nilChan <- 5   // 阻塞
	n := <-nilChan // 阻塞
	fmt.Println(n)

	var closedChan = make(chan int)
	close(closedChan)
	m := <-closedChan
	fmt.Println(m)  // 0
	closedChan <- 5 // panic: send on closed channel
}
