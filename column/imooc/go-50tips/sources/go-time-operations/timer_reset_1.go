package main

import (
	"log"
	"sync"
	"time"
)

func consume(c <-chan bool, timer *time.Timer) bool {
	if !timer.Stop() {
		<-timer.C
	}
	timer.Reset(5 * time.Second)

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
		timer := time.NewTimer(time.Second * 5)
		for {
			if b := consume(c, timer); !b {
				wg.Done()
				return
			}
		}
	}()

	wg.Wait()
}
