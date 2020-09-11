package main

type T struct {
	a int
}

func (t T) M1() {
	t.a = 10
}

func (t *T) M2() {
	t.a = 11
}

func main() {
	var t T // t.a = 0
	println(t.a)

	t.M1()
	println(t.a)

	t.M2()
	println(t.a)
}
