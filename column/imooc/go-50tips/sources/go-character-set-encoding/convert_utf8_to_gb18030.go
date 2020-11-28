package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func dumpToFile(in []byte, filename string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(in)
	if err != nil {
		return err
	}
	return nil
}

func utf8ToGB18030(in []byte) ([]byte, error) {
	if !utf8.Valid(in) {
		return nil, errors.New("invalid utf-8 runes")
	}

	r := bytes.NewReader(in)
	t := transform.NewReader(r, simplifiedchinese.GB18030.NewEncoder())
	out, err := ioutil.ReadAll(t)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func main() {
	var src = "中国人" // <=> "\u4E2D\u56FD\u4EBA"
	var dst []byte

	for i, v := range src {
		fmt.Printf("Unicode字符: %s <=> 码点(rune): %X <=> UTF8编码内存表示: ", string(v), v)
		s := src[i : i+3]
		for _, v := range []byte(s) {
			fmt.Printf("0x%X ", v)
		}

		t, _ := utf8ToGB18030([]byte(s))
		fmt.Printf("<=> GB18030编码内存表示: ")
		for _, v := range t {
			fmt.Printf("0x%X ", v)
		}
		fmt.Printf("\n")

		dst = append(dst, t...)
	}

	dumpToFile(dst, "gb18030.txt")
}
