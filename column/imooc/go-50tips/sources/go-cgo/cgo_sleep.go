package main

//#include <unistd.h>
//void cgoSleep() { sleep(1000); }
import "C"
import (
	"sync"
)

func cgoSleep() {
	C.cgoSleep()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			wg.Done()
			cgoSleep()
		}()
	}

	// 保证所有goroutine都已经启动
	wg.Wait()
	select {}
}
