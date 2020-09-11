package main

import (
	"fmt"
	"io"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func CopyFile(src, dst string) (err error) {
	var r, w *os.File

	// error handler
	defer func() {
		if r != nil {
			r.Close()
		}
		if w != nil {
			w.Close()
		}
		if e := recover(); e != nil {
			if w != nil {
				os.Remove(dst)
			}
			err = fmt.Errorf("copy %s %s: %v", src, dst, err)
		}
	}()

	r, err = os.Open(src)
	check(err)

	w, err = os.Create(dst)
	check(err)

	_, err = io.Copy(w, r)
	check(err)

	return nil
}

func main() {
	err := CopyFile("foo.txt", "bar.txt")
	if err != nil {
		fmt.Println("copyfile error:", err)
		return
	}
	fmt.Println("copyfile ok")
}
