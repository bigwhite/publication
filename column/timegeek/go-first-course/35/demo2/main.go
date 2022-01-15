package main

import (
	"fmt"
	"time"

	"github.com/bigwhite/workerpool"
)

func main() {
	p := workerpool.New(5, workerpool.WithPreAllocWorkers(false), workerpool.WithBlock(false))

	time.Sleep(2 * time.Second)
	for i := 0; i < 10; i++ {
		err := p.Schedule(func() {
			time.Sleep(time.Second * 3)
		})
		if err != nil {
			fmt.Printf("task[%d]: error: %s\n", i, err.Error())
		}
	}

	p.Free()
}
