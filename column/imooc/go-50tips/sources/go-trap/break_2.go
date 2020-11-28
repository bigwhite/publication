package main

import (
	"fmt"
	"time"
)

func breakWithForSwitch(b bool) {
outerloop:
	for {
		time.Sleep(1 * time.Second)
		fmt.Println("enter for-switch loop!")
		switch b {
		case true:
			break outerloop
		case false:
			fmt.Println("go on for-switch loop!")
		}
	}
	fmt.Println("exit breakWithForSwitch")
}

func breakWithForSelect(c <-chan int) {
outerloop:
	for {
		time.Sleep(1 * time.Second)
		fmt.Println("enter for-select loop!")
		select {
		case <-c:
			break outerloop
		default:
			fmt.Println("go on for-select loop!")
		}
	}
	fmt.Println("exit breakWithForSelect")
}

func main() {
	breakWithForSwitch(true)

	c := make(chan int, 1)
	c <- 11
	breakWithForSelect(c)

}
