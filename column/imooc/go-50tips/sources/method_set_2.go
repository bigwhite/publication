package main

type Interface interface {
	M1()
	M2()
}

type T struct{}

func (t T) M1()  {}
func (t *T) M2() {}

func main() {
	var t T
	var pt *T
	DumpMethodSet(&t)
	DumpMethodSet(&pt)
	DumpMethodSet((*Interface)(nil))
}
