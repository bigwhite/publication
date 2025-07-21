package main

import (
	"fmt"
	"os"
	"runtime/debug"
)

func Bar() {
	fmt.Println("In Bar, about to print stack to stderr via debug.PrintStack():")
	debug.PrintStack()
}

func Foo() {
	fmt.Println("In Foo, calling Bar.")
	Bar()
}

func main() {
	fmt.Println("Starting main.")
	Foo()

	if err := someOperationThatMightError(); err != nil {
		// 将堆栈信息作为错误上下文的一部分
		detailedError := fmt.Errorf("operation failed: %w\nCall stack:\n%s", err, debug.Stack())
		fmt.Fprintf(os.Stderr, "%v\n", detailedError)
	}
	fmt.Println("Finished main.")
}

func someOperationThatMightError() error {
	// 模拟一个操作，该操作内部可能还调用了其他函数
	return performComplexStep()
}
func performComplexStep() error {
	return fmt.Errorf("a simulated error occurred deep in call stack")
}
