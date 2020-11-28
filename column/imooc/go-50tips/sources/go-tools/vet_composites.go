package main

import (
	"fmt"
	"net"
)

type myFoo struct {
	name string
	age  int
}

func main() {
	f := &myFoo{"tony", 20}
	err := &net.AddrError{"not found", "localhost"}
	fmt.Printf("%#v,%#v\n", *f, err)
}
