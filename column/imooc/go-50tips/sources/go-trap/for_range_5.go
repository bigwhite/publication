package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var a = []int{1, 2, 3, 4, 5}
	var wg sync.WaitGroup

	for _, v := range a {
		wg.Add(1)
		go func() {
			time.Sleep(time.Second)
			fmt.Println(v)
			wg.Done()
		}()
	}
	wg.Wait()
}
