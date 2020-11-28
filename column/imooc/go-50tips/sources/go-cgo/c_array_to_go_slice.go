package main

// int cArray[] = {1, 2, 3, 4, 5, 6, 7};
import "C"
import (
	"fmt"
	"unsafe"
)

func CArrayToGoArray(cArray unsafe.Pointer, elemSize uintptr, len int) (goArray []int32) {
	for i := 0; i < len; i++ {
		j := *(*int32)((unsafe.Pointer)(uintptr(cArray) + uintptr(i)*elemSize))
		goArray = append(goArray, j)
	}

	return
}

func main() {
	goArray := CArrayToGoArray(unsafe.Pointer(&C.cArray[0]), unsafe.Sizeof(C.cArray[0]), 7)
	fmt.Println(goArray)
}
