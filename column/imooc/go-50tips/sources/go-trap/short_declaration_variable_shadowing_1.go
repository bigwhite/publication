package main

var a int = 13

func main() {
	println(a, &a)
	var a int = 23
	println(a, &a)

	if a == 23 {
		var a int = 33
		println(a, &a)
	}
}
