package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(j int) {
	time.Sleep(time.Second * (time.Duration(j)))
}

func spawnGroup(n int, f func(int)) chan struct{} {
	quit := make(chan struct{})
	job := make(chan int)
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done() // 保证wg.Done在goroutine退出前被执行
			name := fmt.Sprintf("worker-%d:", i)
			for {
				j, ok := <-job
				if !ok {
					println(name, "done")
					return
				}
				// do the job
				worker(j)
			}
		}(i)
	}

	go func() {
		<-quit
		close(job) // 广播给所有新goroutine
		wg.Wait()
		quit <- struct{}{}
	}()

	return quit
}

func main() {
	quit := spawnGroup(5, worker)
	println("spawn a group of workers")

	time.Sleep(5 * time.Second)
	// notify the worker goroutine group to exit
	println("notify the worker group to exit...")
	quit <- struct{}{}

	timer := time.NewTimer(time.Second * 5)
	defer timer.Stop()
	select {
	case <-timer.C:
		println("wait group workers exit timeout!")
	case <-quit:
		println("group workers done")
	}
}
