package main

import (
	"fmt"
	"reflect"
)

func main() {
	// 简单原生类型
	var b = true // 布尔类型
	val := reflect.ValueOf(b)
	typ := reflect.TypeOf(b)
	fmt.Println(typ.Name(), val.Bool()) // bool true

	var i = 23 // 整型
	val = reflect.ValueOf(i)
	typ = reflect.TypeOf(i)
	fmt.Println(typ.Name(), val.Int()) // int 23

	var f = 3.14 // 浮点型
	val = reflect.ValueOf(f)
	typ = reflect.TypeOf(f)
	fmt.Println(typ.Name(), val.Float()) // float64 3.14

	var s = "hello, reflection" // 字符串
	val = reflect.ValueOf(s)
	typ = reflect.TypeOf(s)
	fmt.Println(typ.Name(), val.String()) //string hello, reflection

	var fn = func(a, b int) int { // 函数(一等公民)
		return a + b
	}
	val = reflect.ValueOf(fn)
	typ = reflect.TypeOf(fn)
	fmt.Println(typ.Kind(), typ.String()) // func func(int, int) int

	// 原生复合类型
	var sl = []int{5, 6} // 切片
	val = reflect.ValueOf(sl)
	typ = reflect.TypeOf(sl)
	fmt.Printf("[%d %d]\n", val.Index(0).Int(),
		val.Index(1).Int()) // [5, 6]
	fmt.Println(typ.Kind(), typ.String()) // slice []int

	var arr = [3]int{5, 6} // 数组
	val = reflect.ValueOf(arr)
	typ = reflect.TypeOf(arr)
	fmt.Printf("[%d %d %d]\n", val.Index(0).Int(),
		val.Index(1).Int(), val.Index(2).Int()) // [5 6 0]
	fmt.Println(typ.Kind(), typ.String()) // array [3]int

	var m = map[string]int{ // map
		"tony": 1,
		"jim":  2,
		"john": 3,
	}
	val = reflect.ValueOf(m)
	typ = reflect.TypeOf(m)
	iter := val.MapRange()
	fmt.Printf("{")
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		fmt.Printf("%s:%d,", k.String(), v.Int())
	}
	fmt.Printf("}\n")                     // {tony:1,jim:2,john:3,}
	fmt.Println(typ.Kind(), typ.String()) // map map[string]int

	type Person struct {
		Name string
		Age  int
	}

	var p = Person{"tony", 23} // 结构体
	val = reflect.ValueOf(p)
	typ = reflect.TypeOf(p)
	fmt.Printf("{%s, %d}\n", val.Field(0).String(),
		val.Field(1).Int()) // {"tony", 23}

	fmt.Println(typ.Kind(), typ.Name(), typ.String()) // struct Person main.Person

	var ch = make(chan int, 1) // channel
	val = reflect.ValueOf(ch)
	typ = reflect.TypeOf(ch)
	ch <- 17
	v, ok := val.TryRecv()
	if ok {
		fmt.Println(v.Int()) // 17
	}
	fmt.Println(typ.Kind(), typ.String()) // chan chan int

	// 其他自定义类型
	type MyInt int

	var mi MyInt = 19
	val = reflect.ValueOf(mi)
	typ = reflect.TypeOf(mi)
	fmt.Println(typ.Name(), typ.Kind(), typ.String(), val.Int()) // MyInt int main.MyInt 19
}
