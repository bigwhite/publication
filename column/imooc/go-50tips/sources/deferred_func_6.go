package main

func bar() (int, int) {
	return 1, 2
}

func foo() {
	// builtin functions:
	//	append cap close complex copy delete imag len
	// 	make new panic print println real recover

	var c chan int
	var sl []int
	var m = make(map[string]int, 10)
	m["item1"] = 1
	m["item2"] = 2
	var a = complex(1.0, -1.4)

	var sl1 []int

	defer bar()
	defer append(sl, 11)
	defer cap(sl)
	defer close(c)
	defer complex(2, -2)
	defer copy(sl1, sl)
	defer delete(m, "item2")
	defer imag(a)
	defer len(sl)
	defer make([]int, 10)
	defer new(*int)
	defer panic(1)
	defer print("hello, defer\n")
	defer println("hello, defer")
	defer real(a)
	defer recover()
}

func main() {
	foo()
}
