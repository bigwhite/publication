package main

import "fmt"

type binaryCalcFunc func(int, int) int

func main() {
	var i interface{} = binaryCalcFunc(func(x, y int) int { return x + y })
	c := make(chan func(int, int) int, 10)
	fns := []binaryCalcFunc{
		func(x, y int) int { return x + y },
		func(x, y int) int { return x - y },
		func(x, y int) int { return x * y },
		func(x, y int) int { return x / y },
		func(x, y int) int { return x % y },
	}

	c <- func(x, y int) int {
		return x * y
	}

	fmt.Println(fns[0](5, 6))
	f := <-c
	fmt.Println(f(7, 10))
	v, ok := i.(binaryCalcFunc)
	if !ok {
		fmt.Println("type assertion error")
		return
	}

	fmt.Println(v(17, 7))
}
