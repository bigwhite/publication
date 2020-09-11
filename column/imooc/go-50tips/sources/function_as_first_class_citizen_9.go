package main

import "fmt"

func factorial(n int, f func(int)) {
	if n == 1 {
		f(1) //基本情况
	} else {
		factorial(n-1, func(y int) { f(n * y) })
	}
}

func main() {
	factorial(5, func(y int) { fmt.Printf("%d\n", y) })
}
