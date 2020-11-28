package main

import (
	"fmt"
	"time"
)

func create_timer_by_afterfunc() {
	_ = time.AfterFunc(1*time.Second, func() {
		fmt.Println("timer created by afterfunc fired!")
	})
}

func create_timer_by_newtimer() {
	timer := time.NewTimer(2 * time.Second)
	select {
	case <-timer.C:
		fmt.Println("timer created by newtimer fired!")
	}
}

func create_timer_by_after() {
	select {
	case <-time.After(2 * time.Second):
		fmt.Println("timer created by after fired!")
	}
}

func main() {
	create_timer_by_afterfunc()
	create_timer_by_newtimer()
	create_timer_by_after()
}
