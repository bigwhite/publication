package main

import "fmt"

type foo struct {
	name string
	age  int
}

func (f *foo) setName(name string) {
	f.name = name
}

func (*foo) doSomethingWithoutChange() {
	fmt.Println("doSomethingWithoutChange")
}

type MyInterface interface {
	doSomethingWithoutChange()
}

func main() {
	var strs []string = nil
	strs[0] = "go" // panic

	var m map[string]int
	m["key1"] = 1 // panic

	var f *foo = nil
	f.setName("tony") // panic

	var i MyInterface = nil
	i.doSomethingWithoutChange() // panic
}
