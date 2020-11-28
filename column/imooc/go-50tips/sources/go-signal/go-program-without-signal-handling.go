package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Signal!\n")
	})
	wg.Add(1)
	go func() {
		errChan <- http.ListenAndServe("localhost:8080", nil)
		wg.Done()
	}()

	select {
	case <-time.After(2 * time.Second):
		fmt.Println("web server start ok")
	case err := <-errChan:
		fmt.Println("web server start failed:", err)
	}
	wg.Wait()
	fmt.Println("web server shutdown ok")
}
