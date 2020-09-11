package main

import (
	"errors"
	"fmt"
	"time"
)

var OK = errors.New("ok")

func worker(args ...interface{}) error {
	if len(args) == 0 {
		return errors.New("invalid args")
	}
	interval, ok := args[0].(int)
	if !ok {
		return errors.New("invalid interval arg")
	}

	time.Sleep(time.Second * (time.Duration(interval)))
	return OK
}

func spawn(f func(args ...interface{}) error, args ...interface{}) chan error {
	c := make(chan error)
	go func() {
		c <- f(args...)
	}()
	return c
}

func main() {
	done := spawn(worker, 5)
	println("spawn worker1")
	err := <-done
	fmt.Println("worker1 done:", err)
	done = spawn(worker)
	println("spawn worker2")
	err = <-done
	fmt.Println("worker2 done:", err)
}
