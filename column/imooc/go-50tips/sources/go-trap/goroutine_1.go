package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan error, 1)

	go func() {
		//do something
		time.Sleep(time.Second * 2)
		c <- nil // or c <- errors.New("some error")
	}()

	err := <-c
	fmt.Printf("sub goroutine exit with error: %v\n", err)
}
