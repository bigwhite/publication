package main

import (
	"fmt"
	"reflect"
)

func main() {
	var i = 5
	val := reflect.ValueOf(i)
	r := val.Interface().(int)
	fmt.Println(r) // 5
	r = 6
	fmt.Println(i, r) // 5 6

	val = reflect.ValueOf(&i)
	q := val.Interface().(*int)
	fmt.Printf("%p, %p, %d\n", &i, q, *q) // 0xc0000b4008, 0xc0000b4008, 5
	*q = 7
	fmt.Println(i) // 7
}
