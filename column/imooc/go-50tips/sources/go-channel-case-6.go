package main

import (
	"fmt"
	"time"
)

type counter struct {
	c chan int
	i int
}

var cter counter

func InitCounter() {
	cter = counter{
		c: make(chan int),
	}

	go func() {
		for {
			cter.i++
			cter.c <- cter.i
		}
	}()
	fmt.Println("counter init ok")
}

func Increase() int {
	return <-cter.c
}

func init() {
	InitCounter()
}

func main() {
	for i := 0; i < 10; i++ {
		go func(i int) {
			v := Increase()
			fmt.Printf("goroutine-%d: current counter value is %d\n", i, v)
		}(i)
	}

	time.Sleep(5 * time.Second)
}
