package main

import (
	"sync"
	"time"
)

func goSleep() {
	time.Sleep(time.Second * 1000)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			wg.Done()
			goSleep()
		}()
	}

	//保证所有goroutine都已经启动
	wg.Wait()
	select {}
}
