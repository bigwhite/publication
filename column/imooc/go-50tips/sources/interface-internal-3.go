package main

import "unsafe"

type T int

func (t T) Error() string {
	return "bad error"
}

func main() {
	var eif interface{} = T(5)
	var err error = T(5)
	println("eif:", eif)
	println("err:", err)
	println("eif = err:", eif == err)

	dumpEface(eif)
	dumpItabOfIface(unsafe.Pointer(&err))
	dumpDataOfIface(err)
}
