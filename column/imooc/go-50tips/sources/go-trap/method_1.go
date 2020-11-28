package main

import "fmt"

type foo struct {
	name string
	age  int
}

func (f foo) setNameByValueReceiver(name string) {
	f.name = name
}

func (p *foo) setNameByPointerReceiver(name string) {
	p.name = name
}

func main() {
	f := foo{
		name: "tony",
		age:  20,
	}
	fmt.Println(f) // {tony 20}

	f.setNameByValueReceiver("alex")
	fmt.Println(f) // {tony 20}
	f.setNameByPointerReceiver("alex")
	fmt.Println(f) // {alex 20}
}
