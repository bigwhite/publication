package main

import (
	"fmt"
	"os"
	"sync"
)

func writeToFile(fname string, data []byte, mu *sync.Mutex) error {
	mu.Lock()
	defer mu.Unlock()
	f, err := os.OpenFile(fname, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Seek(0, 2)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return f.Sync()
}

func main() {
	var mu sync.Mutex
	err := writeToFile("./foo.txt", []byte("hello, defer!\n"), &mu)
	if err != nil {
		fmt.Println("writeToFile error:", err)
		return
	}
	fmt.Println("writeToFile ok")
}
