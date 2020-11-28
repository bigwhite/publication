package main

// enum color {
//    RED,
//    BLUE,
//    YELLOW
// };
import "C"
import "fmt"

func main() {
	var e, f, g C.enum_color = C.RED, C.BLUE, C.YELLOW
	fmt.Println(e, f, g) // 0 1 2
}
