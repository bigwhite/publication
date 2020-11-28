package main

import "errors"

func foo() (int, error) {
	return 11, nil
}

func bar() (int, error) {
	return 21, errors.New("error in bar")
}

func main() {
	var err error
	defer func() {
		if err != nil {
			println("error in defer:", err.Error())
		}
	}()

	a, err := foo()
	if err != nil {
		return
	}
	println("a=", a)

	var b int
	if a == 11 {
		b, err = bar()
		if err != nil {
			return
		}
		println("b=", b)
	}
	println("no error occurs")
}
