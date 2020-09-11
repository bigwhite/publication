package main

import "fmt"

func foo1() {
	sl := []int{1, 2, 3}
	defer func(a []int) {
		fmt.Println(a)
	}(sl)

	sl = []int{3, 2, 1}
	_ = sl
}
func foo2() {
	sl := []int{1, 2, 3}
	defer func(p *[]int) {
		fmt.Println(*p)
	}(&sl)

	sl = []int{3, 2, 1}
	_ = sl
}

func main() {
	foo1()
	foo2()
}
