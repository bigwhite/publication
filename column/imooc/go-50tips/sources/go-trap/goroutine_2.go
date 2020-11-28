package main

import "time"

func main() {
	println("main goroutine: start to work...")
	go func() {
		println("goroutine1: start to work...")
		time.Sleep(50 * time.Microsecond)
		println("goroutine1: work done!")
	}()

	go func() {
		println("goroutine2: start to work...")
		time.Sleep(30 * time.Microsecond)
		println("goroutine2: work done!")
	}()

	println("main goroutine: work done!")
}
