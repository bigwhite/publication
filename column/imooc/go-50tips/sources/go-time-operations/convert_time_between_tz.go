package main

import (
	"fmt"
	"time"
)

func main() {
	locSrc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		fmt.Println("load time location failed:", err)
		return
	}

	t := time.Date(2020, 6, 18, 06, 0, 0, 0, locSrc)
	fmt.Println(t) // 美国西部洛杉矶时间，即太平洋时间

	locTo, err := time.LoadLocation("America/New_York")
	if err != nil {
		fmt.Println("load time location failed:", err)
		return
	}

	t1 := t.In(locTo) // 转换成美国东部纽约时间表示
	fmt.Println(t1)
}
