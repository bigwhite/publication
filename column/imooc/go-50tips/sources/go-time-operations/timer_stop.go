package main

import (
	"log"
	"sync"
	"time"
)

func consume(c <-chan bool) bool {
	timer := time.NewTimer(time.Second * 5)
	defer timer.Stop()
	select {
	case b := <-c:
		if b == false {
			log.Printf("recv false, continue")
			return true
		}
		log.Printf("recv true, return")
		return false
	case <-timer.C:
		log.Printf("timer expired")
		return true
	}
}

func main() {
	c := make(chan bool)
	var wg sync.WaitGroup
	wg.Add(2)

	// producer goroutine
	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(time.Second * 1)
			c <- false
		}

		time.Sleep(time.Second * 1)
		c <- true
		wg.Done()
	}()

	go func() {
		for {
			if b := consume(c); !b {
				wg.Done()
				return
			}
		}
	}()

	wg.Wait()
}
