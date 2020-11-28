package main

// #include <stdlib.h>
// #include <stdio.h>
// #include <errno.h>
// int foo(int i) {
//    errno = 0;
//    if (i > 5) {
//        errno = 8;
//        return i - 5;
//    } else {
//        return i;
//    }
//}
import "C"
import "fmt"

func main() {
	i, err := C.foo(C.int(8))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(i)
	}
}
