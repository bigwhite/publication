package main

func main() {
	var a int = 5
	b, c := "hello", 3.1415
	println(a, b, c)

	// var a int = 6
	// b, c := "world", 1.1

	println(&a, a)
	a, d := 55, "hello, go" // ok
	println(&a)
	println(a, d)
}
