package c

import (
	"fmt"

	"bitbucket.org/bigwhite/d"
)

func CallC() {
	fmt.Println("call C: v1.3.0")
	fmt.Println("   --> call D:")
	fmt.Printf("\t")
	d.CallD()
	fmt.Println("   --> call D end")
}
