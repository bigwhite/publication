package main

import "sync/atomic"

func main() {
	var x uint64 = 17
	x = atomic.AddUint64(&x, 1)
	//atomic.AddUint64(&x, 1)
	println(x)
}
