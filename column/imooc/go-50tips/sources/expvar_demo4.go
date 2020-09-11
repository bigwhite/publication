package main

import (
	"expvar"
	"fmt"
	"net/http"
	"time"
)

var customVar *expvar.Int

func init() {
	customVar = expvar.NewInt("customVar")
	customVar.Set(17)
}

func main() {
	http.Handle("/hi", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	}))

	// 模拟业务逻辑
	go func() {
		//... ...
		for {
			customVar.Add(1)
			time.Sleep(time.Second)
		}
	}()

	fmt.Println(http.ListenAndServe("localhost:8080", nil))
}
