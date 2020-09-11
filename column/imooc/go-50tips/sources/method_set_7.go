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

func main() {
	t := T{}
	t.M1()
	t.M2()
}
