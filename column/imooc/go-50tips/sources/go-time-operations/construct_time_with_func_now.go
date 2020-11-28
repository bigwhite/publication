package main

import (
	"fmt"
	"time"
	"unsafe"
)

func dumpWallAndExt(t time.Time) {
	var hasMonotonic int

	// 输出wall字段的值
	pWall := (*uint64)(unsafe.Pointer(&t))
	fmt.Printf("0x%X\n", *pWall)
	if (1<<63)&(*pWall) != 0 {
		hasMonotonic = 1
	}
	fmt.Printf("hasMonotonic = %d\n", hasMonotonic)

	// 输出ext字段的值
	pExt := (*int64)(unsafe.Pointer((uintptr(unsafe.Pointer(&t)) + unsafe.Sizeof(uint64(0)))))
	fmt.Printf("0x%X\n", *pExt)
	fmt.Printf("%d\n", *pExt/86400/365) // 粗略计算距今的年数
}

func main() {
	t := time.Now()
	dumpWallAndExt(t)
}
