package main

type T struct{}

func (T) M1()  {}
func (*T) M2() {}

type Interface interface {
	M1()
	M2()
}

type T1 T
type Interface1 Interface

func main() {
	var t T
	var pt *T
	var t1 T1
	var pt1 *T1

	DumpMethodSet(&t)
	DumpMethodSet(&t1)

	DumpMethodSet(&pt)
	DumpMethodSet(&pt1)

	DumpMethodSet((*Interface)(nil))
	DumpMethodSet((*Interface1)(nil))
}
