package main

import "unsafe"

func main() {
	var x unsafe.Pointer
	var y uintptr
	x = unsafe.Pointer(y)
	_ = x
}
