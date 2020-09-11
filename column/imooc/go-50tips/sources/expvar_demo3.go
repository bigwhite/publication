package main

import (
	"expvar"
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
)

type CustomVar struct {
	value int64
}

func (v *CustomVar) String() string {
	return strconv.FormatInt(atomic.LoadInt64(&v.value), 10)
}

func (v *CustomVar) Add(delta int64) {
	atomic.AddInt64(&v.value, delta)
}

func (v *CustomVar) Set(value int64) {
	atomic.StoreInt64(&v.value, value)
}

func init() {
	customVar := &CustomVar{
		value: 17,
	}
	expvar.Publish("customVar", customVar)

}

func main() {
	http.Handle("/hi", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	}))
	fmt.Println(http.ListenAndServe("localhost:8080", nil))
}
