package main

import (
	"errors"
	"fmt"
)

type X struct {
	i int
}

type Y struct {
	i int
}

type Z struct {
	i int
	s string
	p *[]int
}

func (y *Y) foo(x *X) {
	panic("panic in foo")
}

func (y *Y) bar(x *X) *Y {
	panic("panic in bar")
	return y
}

func (y *Y) baz(x *X, z Z, a, b, c, d, e int16) error {
	panic("panic in baz")
	return errors.New("error in baz")
}

func (y *Y) bam() {
	panic("panic in bam")
}

func main() {
	y := new(Y)
	x := new(X)
	// comment out the ones you don't want to check
	//y.foo(x)
	//y.bar(x)
	z := Z{
		i: 5,
		s: "hello",
		p: &[]int{1, 2, 4},
	}
	err := y.baz(x, z, 11, 12, 13, 14, 15)
	fmt.Println(err)
	//y.bam()
}

//https://stackoverflow.com/questions/36459827/how-to-interpret-go-stacktrace
