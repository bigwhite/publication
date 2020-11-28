package main

//
// void foo() {
// }
//
import "C"

func CallCFunc() {
	C.foo()
}

func foo() {

}

func CallGoFunc() {
	foo()
}
