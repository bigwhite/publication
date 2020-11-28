package main

// char cArray[] = {'a', 'b', 'c', 'd', 'e', 'f', 'g'};
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	goArray := C.GoBytes(unsafe.Pointer(&C.cArray[0]), 7)
	fmt.Printf("%c\n", goArray) // [a b c d e f g]
}
