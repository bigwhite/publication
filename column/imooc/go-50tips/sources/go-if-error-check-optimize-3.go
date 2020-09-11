package main

import (
	"fmt"
	"io"
	"os"
)

type FileCopier struct {
	w   *os.File
	r   *os.File
	err error
}

func (f *FileCopier) open(path string) (*os.File, error) {
	if f.err != nil {
		return nil, f.err
	}

	h, err := os.Open(path)
	if err != nil {
		f.err = err
		return nil, err
	}
	return h, nil
}

func (f *FileCopier) openSrc(path string) {
	if f.err != nil {
		return
	}

	f.r, f.err = f.open(path)
	return
}

func (f *FileCopier) createDst(path string) {
	if f.err != nil {
		return
	}

	f.w, f.err = os.Create(path)
	return
}

func (f *FileCopier) copy() {
	if f.err != nil {
		return
	}

	if _, err := io.Copy(f.w, f.r); err != nil {
		f.err = err
	}
}

func (f *FileCopier) CopyFile(src, dst string) error {
	if f.err != nil {
		return f.err
	}

	defer func() {
		if f.r != nil {
			f.r.Close()
		}
		if f.w != nil {
			f.w.Close()
		}
		if f.err != nil {
			if f.w != nil {
				os.Remove(dst)
			}
		}
	}()

	f.openSrc(src)
	f.createDst(dst)
	f.copy()
	return f.err
}

func main() {
	var fc FileCopier
	err := fc.CopyFile("foo.txt", "bar.txt")
	if err != nil {
		fmt.Println("copy file error:", err)
		return
	}
	fmt.Println("copy file ok")
}
