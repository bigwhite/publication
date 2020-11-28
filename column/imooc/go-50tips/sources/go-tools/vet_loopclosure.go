package main

import "time"

func main() {
	var s = []int{11, 12, 13, 14}
	for i, v := range s {
		go func() {
			println(i)
			println(v)
		}()
	}
	time.Sleep(5 * time.Second)
}
