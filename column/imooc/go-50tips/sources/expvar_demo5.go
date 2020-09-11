package main

import (
	"expvar"
	"fmt"
	"net/http"
	"time"
)

var customVar *expvar.Map

func init() {
	customVar = expvar.NewMap("customVar")

	var field1 expvar.Int
	var field2 expvar.Float
	customVar.Set("field1", &field1)
	customVar.Set("field2", &field2)
}

func main() {
	http.Handle("/hi", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	}))

	// 模拟业务逻辑
	go func() {
		//... ...
		for {
			customVar.Add("field1", 1)
			customVar.AddFloat("field2", 0.001)
			time.Sleep(time.Second)
		}
	}()

	fmt.Println(http.ListenAndServe("localhost:8080", nil))
}
