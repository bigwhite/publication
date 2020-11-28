package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	var wg sync.WaitGroup

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Signal!\n")
	})
	var srv = http.Server{
		Addr: "localhost:8080",
	}

	srv.RegisterOnShutdown(func() {
		// 在一个单独的goroutine中执行
		fmt.Println("clean resources on shutdown...")
		time.Sleep(2 * time.Second)
		fmt.Println("clean resources ok")
		wg.Done()
	})

	wg.Add(2)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
			syscall.SIGHUP)

		<-quit

		timeoutCtx, cf := context.WithTimeout(context.Background(), time.Second*5)
		defer cf()
		var done = make(chan struct{}, 1)
		go func() {
			if err := srv.Shutdown(timeoutCtx); err != nil {
				fmt.Printf("web server shutdown error: %v", err)
			} else {
				fmt.Println("web server shutdown ok")
			}
			done <- struct{}{}
			wg.Done()
		}()

		select {
		case <-timeoutCtx.Done():
			fmt.Println("web server shutdown timeout")
		case <-done:
		}
	}()

	err := srv.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			fmt.Printf("web server start failed: %v\n", err)
			return
		}
	}
	wg.Wait()
	fmt.Println("program exit ok")
}
