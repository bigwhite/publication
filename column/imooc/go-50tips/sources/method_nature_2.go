package main

type T struct {
	a int
}

func (t T) M1() {
}

func (t *T) M2() {
	t.a = 11
}

func main() {
	var t T
	t.M1()
	t.M2() // <=> (&t).M2()

	var pt = &T{}
	pt.M1() // <=> (*pt).M1()
	pt.M2()
}
