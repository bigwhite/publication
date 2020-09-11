package main

//#include <stdio.h>
// void crash() {
//    int *q = NULL;
//    (*q) = 15000;
//    printf("%d\n", *q);
// }
import "C"

import (
	"fmt"
)

func bar() {
	C.crash()
}

func foo() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("recovered from a panic:", e)
		}
	}()
	bar()
}

func main() {
	foo()
	fmt.Println("main exit normally")
}
