package main

import "time"

func worker(j int) {
	time.Sleep(time.Second * (time.Duration(j)))
}

func spawn(f func(int)) chan string {
	quit := make(chan string)
	go func() {
		var job chan int // 模拟job channel
		for {
			select {
			case j := <-job:
				f(j)
			case <-quit:
				quit <- "ok"
			}
		}
	}()
	return quit
}

func main() {
	quit := spawn(worker)
	println("spawn a worker goroutine")

	time.Sleep(5 * time.Second)

	// notify the child goroutine to exit
	println("notify the worker to exit...")
	quit <- "exit"

	timer := time.NewTimer(time.Second * 10)
	defer timer.Stop()
	select {
	case status := <-quit:
		println("worker done:", status)
	case <-timer.C:
		println("wait worker exit timeout")
	}
}
