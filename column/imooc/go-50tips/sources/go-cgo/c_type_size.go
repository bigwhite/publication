package main

// struct employee {
//     char *id;
//     int  age;
// };
import "C"

import (
	"fmt"
)

func main() {
	fmt.Printf("%#v\n", C.sizeof_int)             // 4
	fmt.Printf("%#v\n", C.sizeof_char)            // 1
	fmt.Printf("%#v\n", C.sizeof_struct_employee) // 16
}
