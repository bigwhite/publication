package main

import (
	"expvar"
	"fmt"
	"net/http"
	"time"
)

type CustomVar struct {
	Field1 int64   `json:"field1"`
	Field2 float64 `json:"field2"`
}

var (
	field1 expvar.Int
	field2 expvar.Float
)

func exportStruct() interface{} {
	return CustomVar{
		Field1: field1.Value(),
		Field2: field2.Value(),
	}
}

func init() {
	expvar.Publish("customVar", expvar.Func(exportStruct))
}

func main() {
	http.Handle("/hi", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	}))

	// 模拟业务逻辑
	go func() {
		//... ...
		for {
			field1.Add(1)
			field2.Add(0.001)
			time.Sleep(time.Second)
		}
	}()

	fmt.Println(http.ListenAndServe("localhost:8080", nil))
}
