package main

import (
	"fmt"
	"unsafe"
)

func main() {
	// original string
	var s string = "hello"
	fmt.Println("original string:", s)

	// try to change the original string through unsafe pointer
	modifyString(&s)
	fmt.Println(s)
}

func modifyString(s *string) {
	// 取出第一个8字节的值
	p := (*uintptr)(unsafe.Pointer(s))

	// 获取底层数组的地址
	var array *[5]byte = (*[5]byte)(unsafe.Pointer(*p))

	var len *int = (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(s)) + unsafe.Sizeof((*uintptr)(nil))))

	for i := 0; i < (*len); i++ {
		fmt.Printf("%p => %c\n", &((*array)[i]), (*array)[i])
		p1 := &((*array)[i])
		v := (*p1)
		(*p1) = v + 1 //try to change the character
	}
}
