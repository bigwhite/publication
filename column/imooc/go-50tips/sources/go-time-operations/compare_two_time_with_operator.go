package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now()
	fmt.Println(t) //北京时间

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		fmt.Println("load time location failed:", err)
		return
	}

	t1 := t.In(loc) // 转换成美国东部纽约时间表示
	fmt.Println(t1)
	fmt.Println(t == t1)
}
