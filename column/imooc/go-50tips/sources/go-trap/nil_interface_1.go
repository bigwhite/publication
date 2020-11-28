package main

import (
	"io"
	"strings"
)

type TxtReader struct{}

func (*TxtReader) Read(p []byte) (n int, err error) {
	// ... ...
	return 0, nil
}

func NewTxtReader(path string) io.Reader {
	var r *TxtReader
	if strings.Contains(path, ".txt") {
		r = new(TxtReader)
	}
	return r
}

func main() {
	i := NewTxtReader("/home/tony/test.png")
	if i == nil {
		println("fail to init txt reader")
		return
	}
	println("init txt reader ok")
}
