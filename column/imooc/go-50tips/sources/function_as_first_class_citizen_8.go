package main

import "fmt"

func Max(n int, m int, f func(y int)) {
	if n > m {
		f(n)
	} else {
		f(m)
	}

}

func main() {
	Max(5, 6, func(y int) { fmt.Printf("%d\n", y) })
}
