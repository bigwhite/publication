package main

import (
	"fmt"
	"log"
	"os"
)

func doSomething() error {
	// ... 模拟操作 ...
	return fmt.Errorf("simulated error: something went wrong")
}

func main() {
	fmt.Println("Application starting...") // 直接输出到 Stdout

	// 配置标准库 log
	log.SetOutput(os.Stdout) // 通常默认是 Stderr，这里改为 Stdout 方便观察
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.Println("Standard logger configured.")

	err := doSomething()
	if err != nil {
		// 不同的输出方式
		fmt.Printf("Output with fmt.Printf: Error occurred: %v\n", err)
		log.Printf("Output with log.Printf: Error occurred: %v\n", err)
	}
	fmt.Println("Application finished.")
}
