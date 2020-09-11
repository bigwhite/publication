package main

type Interface interface {
	M1()
	M2()
}

type T struct {
	Interface
}

func (T) M1() {
	println("T's M1")
}

type S struct{}

func (S) M1() {
	println("S's M1")
}
func (S) M2() {
	println("S's M2")
}

func main() {
	var t = T{
		Interface: S{},
	}

	t.M1()
	t.M2()
}
