package main

type Interface interface {
	M1()
	M2()
}

type T struct {
	Interface
}

func (T) M3() {}

func main() {
	DumpMethodSet((*Interface)(nil))
	var t T
	var pt *T
	DumpMethodSet(&t)
	DumpMethodSet(&pt)
}
