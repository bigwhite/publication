package main

import (
	"demo/mypackage"
	"fmt"
)

func main() {
	ms := mypackage.NewMyStruct("Hello")
	fmt.Println(ms.Field) // 可以访问 Field
	ms.M1()               // 可以调用 M1
}
