package main

import "fmt"

func main() {
	var n int = 61
	var ei interface{} = n
	n = 62
	fmt.Println("data in box:", ei)

	var m int = 51
	ei = &m
	m = 52

	p := ei.(*int)
	fmt.Println("data in box:", *p)
}
