package main

import "fmt"

type foo struct {
	name string
	age  int
}

func (*foo) doSomethingWithoutChange() {
	fmt.Println("doSomethingWithoutChange")
}

type MyInterface interface {
	doSomethingWithoutChange()
}

func main() {
	// 切片仅在调用append操作时才是零值可用的
	var strs []string = nil
	strs = append(strs, "hello", "go")
	fmt.Printf("%q\n", strs)

	// 自定义类型的方法中没有对自身实例进行解引用操作时，
	// 我们可以通过该类型的零值指针调用其方法
	var f *foo = nil
	f.doSomethingWithoutChange()

	// 给接口类型赋予显式转型后的nil(并非真正的零值)
	// 我们可以通过该接口调用没有解引用操作的方法
	var i MyInterface = (*foo)(nil)
	i.doSomethingWithoutChange()
}
