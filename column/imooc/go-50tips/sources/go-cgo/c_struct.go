package main

// #include <stdlib.h>
//
// struct employee {
//     char *id;
//     int  age;
// };
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {
	id := C.CString("1247")
	defer C.free(unsafe.Pointer(id))

	var p = C.struct_employee{
		id:  id,
		age: 21,
	}
	fmt.Printf("%#v\n", p)
}
