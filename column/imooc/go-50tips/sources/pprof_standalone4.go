package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	go func() {
		fmt.Println(http.ListenAndServe("localhost:8080", nil))
	}()

	var wg sync.WaitGroup
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	wg.Add(1)
	go func() {
		for {
			select {
			case <-c:
				wg.Done()
				return
			default:
				s1 := "hello,"
				s2 := "gopher"
				s3 := "!"
				_ = s1 + s2 + s3
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()
	wg.Wait()
	fmt.Println("program exit")
}
