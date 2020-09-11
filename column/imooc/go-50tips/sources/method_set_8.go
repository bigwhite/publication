package main

type Interface interface {
	M1()
	M2()
	M3()
}

type Interface1 interface {
	M1()
	M2()
	M4()
}

type T struct {
	Interface
	Interface1
}

func (T) M1() { println("T's M1") }
func (T) M2() { println("T's M2") }

func main() {
	t := T{}
	t.M1()
	t.M2()
}
