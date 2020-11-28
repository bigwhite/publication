package main

// #include <stdio.h>
// union bar {
//        char   c;
//        int    i;
//        double d;
// };
import "C"
import "fmt"

func main() {
	var b *C.union_bar = new(C.union_bar)
	b.c = 4
	fmt.Println(b)
}
