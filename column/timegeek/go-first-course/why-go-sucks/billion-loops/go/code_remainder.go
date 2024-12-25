package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

func main() {
	input, e := strconv.Atoi(os.Args[1])
	if e != nil {
		panic(e)
	}
	u := int32(input)
	mask := u - 1 // 只在u为2的幂时有效
	r := int32(rand.Intn(10000))
	var a [10000]int32

	for i := int32(0); i < 10000; i++ {
		for j := int32(0); j < 100000; j++ {
			a[i] = a[i] + (j & mask) // 使用位运算代替模运算
		}
		a[i] += r
	}
	fmt.Println(a[r])
}
