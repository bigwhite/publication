package main

import (
	"fmt"
	"reflect"
)

func Add(i, j int) int {
	return i + j
}

type Calculator struct{}

func (c Calculator) Add(i, j int) int {
	return i + j
}

func main() {
	// 函数调用
	f := reflect.ValueOf(Add)
	var i = 5
	var j = 6
	vals := []reflect.Value{reflect.ValueOf(i), reflect.ValueOf(j)}
	ret := f.Call(vals)
	fmt.Println(ret[0].Int()) // 11

	// 方法调用
	c := reflect.ValueOf(Calculator{})
	m := c.MethodByName("Add")
	ret = m.Call(vals)
	fmt.Println(ret[0].Int()) // 11

	var k float64 = 3.14
	ret = m.Call([]reflect.Value{reflect.ValueOf(i),
		reflect.ValueOf(k)}) // panic: reflect: Call using float64 as type int
	fmt.Println(ret[0].Int())
}
