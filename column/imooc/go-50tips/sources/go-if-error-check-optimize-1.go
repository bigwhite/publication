package main

import (
	"fmt"
	"io"
	"os"
)

func openBoth(src, dst string) (*os.File, *os.File, error) {
	var r, w *os.File
	var err error
	if r, err = os.Open(src); err != nil {
		return nil, nil, fmt.Errorf("copy %s %s: %v", src, dst, err)
	}

	if w, err = os.Create(dst); err != nil {
		r.Close()
		return nil, nil, fmt.Errorf("copy %s %s: %v", src, dst, err)
	}
	return r, w, nil
}

func CopyFile(src, dst string) error {
	var err error
	var r, w *os.File
	if r, w, err = openBoth(src, dst); err != nil {
		return err
	}
	defer func() {
		r.Close()
		w.Close()
		if err != nil {
			os.Remove(dst)
		}
	}()

	if _, err = io.Copy(w, r); err != nil {
		return fmt.Errorf("copy %s %s: %v", src, dst, err)
	}
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
